#include "config.hpp"

#include <pistache/net.h>
#include <nlohmann/json.hpp>

#include <fstream>

using namespace Pistache;
using namespace nlohmann;

Config::Config(std::string pathToConfig) {
    std::ifstream configFile(pathToConfig);
    json config = json::parse(configFile);
    std::string port_envar = config["port_envar"].get<std::string>();
    char* port_c_str = std::getenv(port_envar.c_str());
    if (port_c_str == nullptr) {
        port = config["default_port"].get<uint16_t>();
    }
    else {
        port = std::stoi(port_c_str);   
    }
    threads = config["threads"].get<unsigned int>();
    imageUrls = config["image_urls"].get<std::vector<std::string>>();
}
