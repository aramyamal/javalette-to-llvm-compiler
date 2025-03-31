grammar Javalette;

// PARSER RULES
// a program is a list of definitions
prgm
    : def* 
    ;

// definitions have type, identifier, arguments and statements
def 
    : type Ident '(' (arg (',' arg)*)? ')' '{' stm* '}' # FuncDef
    ;

// an argument is a type and identifier
arg
    : type Ident                                # ParamArg
    ;

// statements can be the following
stm
    : exp ';'                                   # ExpStm
    | type item (',' item)* ';'                 # DeclsStm
    | 'return' exp ';'                          # ReturnStm
    | 'return' ';'                              # VoidReturnStm
    | 'while' '(' exp ')' stm                   # WhileStm
    | '{' stm* '}'                              # BlockStm
    | 'if' '(' exp ')' stm ('else' stm)?        # IfStm
    | ';'                                       # BlankStm
    ;

item
    : Ident                             # NoInitItem
    | Ident '=' exp                     # InitItem
    ;

// expressions can be the following
exp
    : '(' exp ')'                       # ParenExp
    | boolLit                           # BoolExp
    | Integer                           # IntExp
    | Double                            # DoubleExp
    | Ident                             # IdentExp
    | Ident '(' (exp (',' exp)*)? ')'   # FuncExp
    | String                            # StringExp
    | '-' exp                           # NegExp
    | '!' exp                           # NotExp
    | Ident incDecOp                    # PostExp
    | incDecOp Ident                    # PreExp
    | exp mulOp exp                     # MulExp
    | exp addOp exp                     # AddExp
    | exp cmpOp exp                     # CmpExp
    | exp '&&' exp                      # AndExp
    | exp '||' exp                      # OrExp
    | <assoc=right> Ident '=' exp       # AssignExp
    ;

boolType: 'boolean';
intType: 'int';
doubleType: 'double';
stringType: 'type';
voidType: 'void';
type: boolType | intType | doubleType | stringType | voidType ;

boolLit
    : 'true'                            #TrueLit
    | 'false'                           #FalseLit
    ;

incDecOp
    : '++'                              #Inc
    | '--'                              #Dec
    ;

mulOp
    : '*'                               #Mul
    | '/'                               #Div
    | '%'                               #Mod
    ;

addOp
    : '+'                               #Add
    | '-'                               #Sub
    ;

cmpOp
    : '<'                               #LTh
    | '>'                               #GTh
    | '<='                              #LTE 
    | '>='                              #GTE
    | '=='                              #Equ
    | '!='                              #NEq
    ;

// LEXER RULES
Ident: Letter (Letter | Digit | '_')*;
Integer: Digit+;
Double: Digit+ '.' Digit+ | Digit+ ('.' Digit+)? ('e' | 'E') ('+' | '-')? Digit+;

String: '"' (~["\\] | '\\' .)* '"';

fragment Letter: [a-zA-Z];
fragment Digit: [0-9];

// skip whitespace and comments
WS: [ \t\r\n]+ -> skip;
LineComment: ('//' | '#') ~[\r\n]* -> skip;
BlockComment: '/*' .*? '*/' -> skip;

