# Javalette to LLVM Compiler
> **Note:**  
> This project was developed as the sole assignment for a master's level university course in compiler construction, where the task was to implement a compiler from scratch.
> 

This project is a compiler for the [Javalette](https://github.com/TDA283-compiler-construction/project/blob/master/project/javalette.md) language (a combination of a subset of C and a subset of Java), targeting LLVM IR. Heap allocated multi-dimensional arrays and heap allocated structs are also implemented according to [these specifications](https://github.com/TDA283-compiler-construction/project/blob/master/project/extensions.md). The compiler is implemented in Go and uses ANTLR for parser generation.

## Building

To build the project, simply run:

```sh
make
```

This will:
- Clean previous builds
- Download ANTLR and generate the parser
- Build the executables

## Usage

After building, you will find two executables in the repository root under `build/`: `jlc` and `typecheck`.

### Compile and Typecheck

- **From Standard Input:**
  ```sh
  ./jlc
  ```
  Paste your Javalette code and press `Ctrl+D` to generate LLVM IR or alternatively pipe the code into the executable.

- **From File:**
  ```sh
  ./jlc <input-file>
  ```

- **With Output File:**
  ```sh
  ./jlc -o <output-file> <input-file>
  ```

### Typecheck Only

- **From File:**
  ```sh
  ./typecheck <input-file>
  ```
