#include "bridge.h"
#include "core/evaluator.h"
#include "core/session.h"
#include "core/numberformatter.h"
#include "core/settings.h"
#include "core/functions.h"
#include "core/constants.h"
#include <QCoreApplication>
#include <QString>
#include <QJsonDocument>
#include <QJsonObject>
#include <QFile>
#include <cstring>
#include <iostream>

static int dummy_argc = 1;
static char* dummy_argv[] = {(char*)"turbocrunch", (char*)NULL};
static QCoreApplication* app = nullptr;
static Session* session = nullptr;

#include <cstdio>

#include "math/units.h"

static char* evaluator_evaluate_internal(EvaluatorPtr p, const char* expression, bool updateAns);

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

char* evaluator_evaluate(EvaluatorPtr p, const char* expression) {
    return evaluator_evaluate_internal(p, expression, false);
}

char* evaluator_evaluate_update_ans(EvaluatorPtr p, const char* expression) {
    return evaluator_evaluate_internal(p, expression, true);
}

static char* evaluator_evaluate_internal(EvaluatorPtr p, const char* expression, bool updateAns) {
    Evaluator* e = (Evaluator*)p;
    QString qexpr = QString::fromUtf8(expression);
    e->setExpression(qexpr);
    
    Quantity res;
    if (updateAns) {
        res = e->evalUpdateAns();
    } else {
        res = e->eval();
    }
    
    QString err = e->error();
    
    QString resultStr;
    if (err.isEmpty()) {
        resultStr = NumberFormatter::format(res);
        resultStr.replace(QChar(0x2212), '-');
        if (resultStr.isEmpty()) {
             resultStr = "Empty Result";
        }
    } else {
        resultStr = QString("Error: ") + err;
    }

    QByteArray bytes = resultStr.toUtf8();
    char* result_ptr = (char*)malloc(bytes.size() + 1);
    if (result_ptr) {
        memcpy(result_ptr, bytes.constData(), bytes.size());
        result_ptr[bytes.size()] = '\0';
    }
    return result_ptr;
}

void evaluator_set_angle_mode(char mode) {
    Settings::instance()->angleUnit = mode;
    Evaluator::instance()->initializeAngleUnits();
}

char evaluator_get_angle_mode() {
    return Settings::instance()->angleUnit;
}

char* evaluator_get_variable(EvaluatorPtr p, const char* name) {
    Evaluator* e = (Evaluator*)p;
    QString qname = QString::fromUtf8(name);
    if (e->hasVariable(qname)) {
        Variable v = e->getVariable(qname);
        QString val = NumberFormatter::format(v.value());
        val.replace(QChar(0x2212), '-');
        
        QByteArray bytes = val.toUtf8();
        char* result_ptr = (char*)malloc(bytes.size() + 1);
        if (result_ptr) {
            memcpy(result_ptr, bytes.constData(), bytes.size());
            result_ptr[bytes.size()] = '\0';
        }
        return result_ptr;
    }
    return nullptr;
}

int evaluator_get_constants_count() {
    return Constants::instance()->list().count();
}

static char* qstring_to_char(const QString& s) {
    QByteArray bytes = s.toUtf8();
    char* ptr = (char*)malloc(bytes.size() + 1);
    if (ptr) {
        memcpy(ptr, bytes.constData(), bytes.size());
        ptr[bytes.size()] = '\0';
    }
    return ptr;
}

char* evaluator_get_constant_name(int index) {
    const QList<Constant>& list = Constants::instance()->list();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].name);
    }
    return nullptr;
}

char* evaluator_get_constant_value(int index) {
    const QList<Constant>& list = Constants::instance()->list();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].value);
    }
    return nullptr;
}

char* evaluator_get_constant_category(int index) {
    const QList<Constant>& list = Constants::instance()->list();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].category);
    }
    return nullptr;
}

int evaluator_get_units_count() {
    return Units::getList().count();
}

char* evaluator_get_unit_name(int index) {
    const QList<Unit>& list = Units::getList();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].name);
    }
    return nullptr;
}

int evaluator_get_functions_count() {
    return FunctionRepo::instance()->getIdentifiers().count();
}

char* evaluator_get_function_identifier(int index) {
    QStringList identifiers = FunctionRepo::instance()->getIdentifiers();
    if (index >= 0 && index < identifiers.count()) {
        return qstring_to_char(identifiers[index]);
    }
    return nullptr;
}

char* evaluator_get_function_name(int index) {
    QStringList identifiers = FunctionRepo::instance()->getIdentifiers();
    if (index >= 0 && index < identifiers.count()) {
        Function* f = FunctionRepo::instance()->find(identifiers[index]);
        if (f) {
            return qstring_to_char(f->name());
        }
    }
    return nullptr;
}

char* evaluator_get_function_usage(int index) {
    QStringList identifiers = FunctionRepo::instance()->getIdentifiers();
    if (index >= 0 && index < identifiers.count()) {
        Function* f = FunctionRepo::instance()->find(identifiers[index]);
        if (f) {
            return qstring_to_char(f->usage());
        }
    }
    return nullptr;
}

int evaluator_get_user_functions_count(EvaluatorPtr p) {
    Evaluator* e = (Evaluator*)p;
    return e->getUserFunctions().count();
}

char* evaluator_get_user_function_name(EvaluatorPtr p, int index) {
    Evaluator* e = (Evaluator*)p;
    QList<UserFunction> list = e->getUserFunctions();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].name());
    }
    return nullptr;
}

char* evaluator_get_user_function_args(EvaluatorPtr p, int index) {
    Evaluator* e = (Evaluator*)p;
    QList<UserFunction> list = e->getUserFunctions();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].arguments().join(";"));
    }
    return nullptr;
}

char* evaluator_get_user_function_expression(EvaluatorPtr p, int index) {
    Evaluator* e = (Evaluator*)p;
    QList<UserFunction> list = e->getUserFunctions();
    if (index >= 0 && index < list.count()) {
        return qstring_to_char(list[index].expression());
    }
    return nullptr;
}

void evaluator_unset_user_function(EvaluatorPtr p, const char* name) {
    Evaluator* e = (Evaluator*)p;
    e->unsetUserFunction(QString::fromUtf8(name));
}

int evaluator_is_user_function_assign(EvaluatorPtr p) {
    Evaluator* e = (Evaluator*)p;
    return e->isUserFunctionAssign() ? 1 : 0;
}

int evaluator_session_save(EvaluatorPtr p, const char* filename) {
    Evaluator* e = (Evaluator*)p;
    Session* s = (Session*)e->session();
    if (!s) return 0;

    QFile file(QString::fromUtf8(filename));
    if (!file.open(QIODevice::WriteOnly)) {
        return 0;
    }

    QJsonObject json;
    s->serialize(json);
    QJsonDocument doc(json);
    file.write(doc.toJson());
    file.close();
    return 1;
}

int evaluator_session_load(EvaluatorPtr p, const char* filename) {
    Evaluator* e = (Evaluator*)p;
    Session* s = (Session*)e->session();
    if (!s) {
        s = new Session();
        e->setSession(s);
    }

    QFile file(QString::fromUtf8(filename));
    if (!file.open(QIODevice::ReadOnly)) {
        return 0;
    }

    QByteArray data = file.readAll();
    QJsonDocument doc(QJsonDocument::fromJson(data));
    s->deSerialize(doc.object(), false);
    file.close();
    
    // After loading, we might need to re-initialize built-in variables
    e->initializeBuiltInVariables();
    
    return 1;
}
