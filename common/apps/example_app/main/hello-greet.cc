#include "common/apps/example_app/main/hello-greet.h"
#include <string>

std::string get_greet(const std::string& who) {
    return "Hello " + who;
}
