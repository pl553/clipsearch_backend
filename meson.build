project('clipsearch_backend', 'cpp',
  version : '0.1',
  default_options : ['warning_level=3',
                     'cpp_std=c++17'])

pistache = dependency('libpistache')

exe = executable('clipsearch_api_server', 'src/main.cpp', 'src/clipsearch_api.cpp', 'src/config.cpp',
  dependencies : [pistache])
