#ifndef BRIDGE_H
#define BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

typedef void* EvaluatorPtr;

EvaluatorPtr evaluator_init();
void evaluator_destroy(EvaluatorPtr p);

/**
 * evaluator_evaluate returns the result of the expression.
 * The caller is responsible for calling free() on the returned string.
 * If there is an error, it returns a string starting with "Error: ".
 */
char* evaluator_evaluate(EvaluatorPtr p, const char* expression);

#ifdef __cplusplus
}
#endif

#endif
