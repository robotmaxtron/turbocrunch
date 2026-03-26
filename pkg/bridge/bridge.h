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
char* evaluator_evaluate_update_ans(EvaluatorPtr p, const char* expression);

void evaluator_set_angle_mode(char mode); // 'r', 'd', 'g'
char evaluator_get_angle_mode();

char* evaluator_get_variable(EvaluatorPtr p, const char* name);

// Constants
int evaluator_get_constants_count();
char* evaluator_get_constant_name(int index);
char* evaluator_get_constant_value(int index);
char* evaluator_get_constant_category(int index);

// Units
int evaluator_get_units_count();
char* evaluator_get_unit_name(int index);

// Functions
int evaluator_get_functions_count();
char* evaluator_get_function_identifier(int index);
char* evaluator_get_function_name(int index);
char* evaluator_get_function_usage(int index);

// User-Defined Functions
int evaluator_get_user_functions_count(EvaluatorPtr p);
char* evaluator_get_user_function_name(EvaluatorPtr p, int index);
char* evaluator_get_user_function_args(EvaluatorPtr p, int index);
char* evaluator_get_user_function_expression(EvaluatorPtr p, int index);
void evaluator_unset_user_function(EvaluatorPtr p, const char* name);
int evaluator_is_user_function_assign(EvaluatorPtr p);
int evaluator_session_save(EvaluatorPtr p, const char* filename);
int evaluator_session_load(EvaluatorPtr p, const char* filename);

// Variables
int evaluator_get_variables_count(EvaluatorPtr p);
char* evaluator_get_variable_name(EvaluatorPtr p, int index);
char* evaluator_get_variable_value(EvaluatorPtr p, int index);
void evaluator_unset_variable(EvaluatorPtr p, const char* name);

void evaluator_set_result_format(char format); // 'd', 'h', 'b', 'o'
char evaluator_get_result_format();

#ifdef __cplusplus
}
#endif

#endif
