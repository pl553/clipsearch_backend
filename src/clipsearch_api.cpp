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
    j["data"]["image_urls"] = _config.imageUrls;
    response.send(Http::Code::Ok, j.dump(), Http::Mime::MediaType::fromString("application/json"));
}
