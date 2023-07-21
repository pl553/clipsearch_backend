#pragma once

#include <pistache/net.h>

#include <string>
#include <vector>

class Config {
public:
    explicit Config(std::string pathToConfig);
    Pistache::Port port;
    unsigned int threads;
    std::vector<std::string> imageUrls;
};
