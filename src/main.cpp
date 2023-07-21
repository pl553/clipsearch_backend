#include "clipsearch_api.hpp"
#include "config.hpp"

#include <pistache/endpoint.h>
#include <pistache/router.h>
#include <pistache/http.h>

#include <signal.h>

#include <algorithm>
#include <iostream>
#include <vector>
#include <string>

using namespace Pistache;
using namespace Pistache::Rest;

void waitForShutdownRequest()
{
  sigset_t sigset;

  sigemptyset(&sigset);
  sigaddset(&sigset, SIGHUP);
  sigaddset(&sigset, SIGINT);
  sigaddset(&sigset, SIGTERM);
  sigaddset(&sigset, SIGQUIT);
  sigprocmask(SIG_BLOCK, &sigset, nullptr);

  int sig = 0;
  sigwait(&sigset, &sig);
  std::cerr << "\nReceived signal: " << sig << ", " << strsignal(sig) << "\n";

  sigprocmask(SIG_UNBLOCK, &sigset, nullptr);
}

int main(int argc, char *argv[]) {
    std::vector<std::string> args(argv, argv + argc);
    std::string configPath = "config.json";
    auto it = std::find(args.begin(), args.end(), "-c");
    if (it != args.end() && ++it != args.end()) {
        configPath = *it;
    }

    Config config(configPath);
    ClipSearchApiController apiController(config);
    
    Address addr(Ipv4::any(), config.port);

    auto opts = Http::Endpoint::options().threads(config.threads);
    Http::Endpoint server(addr);
    server.init(opts);
    
    Router router;
    Routes::Get(router, "/api/gallery", Routes::bind(&ClipSearchApiController::OnGalleryRequest, &apiController)); 
    
    server.setHandler(router.handler());
    server.serveThreaded();
    
    waitForShutdownRequest();
    
    server.shutdown();
}
