#include "config.hpp"

#include <pistache/net.h>
#include <nlohmann/json.hpp>

#include <fstream>

using namespace Pistache;
using namespace nlohmann;

Config::Config(std::string pathToConfig) {
    std::ifstream configFile(pathToConfig);
    json config = json::parse(configFile);
    port = Port(config["port"].get<uint16_t>());
    threads = config["threads"].get<unsigned int>();
    imageUrls = config["image_urls"].get<std::vector<std::string>>();
}
