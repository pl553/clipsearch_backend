#include "clipsearch_api.hpp"

#include <pistache/endpoint.h>
#include <pistache/router.h>
#include <pistache/http.h>
#include <nlohmann/json.hpp>

#include <iostream>

using namespace Pistache;
using namespace nlohmann;

ClipSearchApiController::ClipSearchApiController(Config config) : _config(config) {
}

void ClipSearchApiController::OnGalleryRequest(const Rest::Request& request, Http::ResponseWriter response) {
    json j;
    j["status"] = "success";
    
    int mask = (~_requestCount) & 7;
    j["data"]["image_urls"] = json::array();
    for (int i = 0; i < 3; ++i) {
        if (mask & (1 << i)) {
            j["data"]["image_urls"].push_back(_config.imageUrls[i]);
        }
    }
    
    ++_requestCount;
    
    response.send(Http::Code::Ok, j.dump(), Http::Mime::MediaType::fromString("application/json"));
}
