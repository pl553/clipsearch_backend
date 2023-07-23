#pragma once

#include "config.hpp"

#include <pistache/endpoint.h>
#include <pistache/router.h>
#include <pistache/http.h>

#include <atomic>

using namespace Pistache;

class ClipSearchApiController {
public:
    explicit ClipSearchApiController(Config config);
    void OnGalleryRequest(const Rest::Request& request, Http::ResponseWriter response);
private:
    Config _config;
    std::atomic<int> _requestCount = 0;
};
