#include "bridge.h"
#include "core/evaluator.h"
#include "core/numberformatter.h"
#include "core/settings.h"
#include "core/functions.h"
#include "core/constants.h"
#include <QCoreApplication>
#include <QString>
#include <cstring>
#include <iostream>

static int dummy_argc = 1;
static char* dummy_argv[] = {(char*)"turbocrunch", (char*)NULL};
static QCoreApplication* app = nullptr;

#include <cstdio>

#include "math/units.h"

static void initialize_speedcrunch() {
    if (!QCoreApplication::instance()) {
        app = new QCoreApplication(dummy_argc, dummy_argv);
    }
    
    // Initialize singletons and repos
    Settings::instance()->load();
    Settings::instance()->complexNumbers = true;
    
    FunctionRepo::instance();
    Constants::instance();
    
    // Force units initialization
    Units::getList();
    
    Evaluator* e = Evaluator::instance();
    e->reset(); // This calls initializeBuiltInVariables()
    e->initializeAngleUnits();
}

EvaluatorPtr evaluator_init() {
    static bool initialized = false;
    if (!initialized) {
        initialize_speedcrunch();
        initialized = true;
    }
    return (EvaluatorPtr)Evaluator::instance();
}

void evaluator_destroy(EvaluatorPtr p) {
    // Evaluator is a singleton in SpeedCrunch
}

const char* evaluator_evaluate(EvaluatorPtr p, const char* expression) {
    Evaluator* e = (Evaluator*)p;
    QString qexpr = QString::fromUtf8(expression);
    e->setExpression(qexpr);
    
    Quantity res = e->eval();
    
    QString err = e->error();
    
    static char result_buf[4096];
    if (err.isEmpty()) {
        QString formatted = NumberFormatter::format(res);
        formatted.replace(QChar(0x2212), '-');
        if (formatted.isEmpty()) {
             strncpy(result_buf, "Empty Result", sizeof(result_buf)-1);
        } else {
             strncpy(result_buf, formatted.toUtf8().constData(), sizeof(result_buf)-1);
        }
    } else {
        strncpy(result_buf, (QString("Error: ") + err).toUtf8().constData(), sizeof(result_buf)-1);
    }
    result_buf[sizeof(result_buf)-1] = '\0';
    return result_buf;
}
