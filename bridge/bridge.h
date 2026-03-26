#ifndef BRIDGE_H
#define BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef void* EvaluatorPtr;

EvaluatorPtr evaluator_init();
void evaluator_destroy(EvaluatorPtr p);
const char* evaluator_evaluate(EvaluatorPtr p, const char* expression);

#ifdef __cplusplus
}
#endif

#endif
