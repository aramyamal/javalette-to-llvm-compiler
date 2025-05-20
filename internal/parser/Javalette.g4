grammar Javalette;

// PARSER RULES
// a program is a list of definitions
prgm
    : def* 
    ;

// defintions can be function defs, struct defs and typedef defs
def 
    : type Ident '(' (arg (',' arg)*)? ')' '{' stm* '}' # FuncDef
    | 'struct' Ident '{' structField* '}' ';'           # StructDef
    | 'typedef' 'struct' Ident '*' Ident ';'            # TypedefDef
    ;

// an argument is a type and identifier
arg
    : type Ident                                # ParamArg
    ;

structField
    : type Ident ';'
    ;

// statements can be the following
stm
    : exp ';'                                   # ExpStm
    | type item (',' item)* ';'                 # DeclsStm
    | 'return' exp ';'                          # ReturnStm
    | 'return' ';'                              # VoidReturnStm
    | 'for' '(' type Ident ':' exp ')' stm      # ForEachStm
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
    : '(' exp ')'                                # ParenExp
    | '(' type ')' 'null'                        # NullPtrExp
    | boolLit                                    # BoolExp
    | Integer                                    # IntExp
    | Double                                     # DoubleExp
    | 'new' baseType arrayIndex+                 # NewArrExp
    | 'new' Ident                                # NewStructExp
    | Ident                                      # IdentExp
    | Ident '(' (exp (',' exp)*)? ')'            # FuncExp
    | exp arrayIndex+                            # ArrIndexExp
    | exp '.' Ident                              # FieldExp
    | exp '->' Ident                             # DerefExp
    | String                                     # StringExp
    | '-' exp                                    # NegExp
    | '!' exp                                    # NotExp
    | exp incDecOp                               # PostExp
    | incDecOp exp                               # PreExp
    | exp mulOp exp                              # MulExp
    | exp addOp exp                              # AddExp
    | exp cmpOp exp                              # CmpExp
    | exp '&&' exp                               # AndExp
    | exp '||' exp                               # OrExp
    | <assoc=right> exp '=' exp                  # AssignExp
    ;

arrayIndex
    : '[' exp ']'
    ;

boolType: 'boolean';
intType: 'int';
doubleType: 'double';
stringType: 'string';
voidType: 'void';
baseType
    : boolType
    | intType
    | doubleType
    | stringType
    | voidType
    ;

type
    : baseType arraySuffix*             #PrimitiveType
    | Ident arraySuffix*                #CustomType
    ;

arraySuffix
    : '[' ']'
    ;


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

