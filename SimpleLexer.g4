/*
 * Lexer для упрощенного подмножества Go
 * Поддерживает: переменные, базовые типы, операции, if-else, for, fmt.Println
 */

lexer grammar SimpleLexer;

// Keywords
VAR         : 'var';
TYPE        : 'type';
INT         : 'int';
STRING      : 'string';
BOOL        : 'bool';
IF          : 'if';
ELSE        : 'else';
FOR         : 'for';
FUNC        : 'func';
PRINTLN     : 'fmt.Println';
TRUE        : 'true';
FALSE       : 'false';
STRUCT      : 'struct';

// Identifiers
IDENTIFIER  : [a-zA-Z_][a-zA-Z0-9_]*;

// Literals
INT_LIT     : [0-9]+;
STRING_LIT  : '"' (~["\\\n])* '"';

// Operators
PLUS        : '+';
MINUS       : '-';
STAR        : '*';
DIV         : '/';
ASSIGN      : '=';
DECLARE     : ':=';
EQ          : '==';
NEQ         : '!=';
LT          : '<';
GT          : '>';
LE          : '<=';
GE          : '>=';
AND         : '&&';
OR          : '||';
NOT         : '!';

// Punctuation
LPAREN      : '(';
RPAREN      : ')';
LBRACE      : '{';
RBRACE      : '}';
LBRACKET    : '[';
RBRACKET    : ']';
SEMICOLON   : ';';
COMMA       : ',';
COLON       : ':';
DOT         : '.';

// Whitespace and comments
WS          : [ \t\n\r]+ -> skip;
COMMENT     : '//' ~[\n]* -> skip;
MCOMMENT    : '/*' .*? '*/' -> skip;
