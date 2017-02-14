# Plis
Plis is an application that makes it easier for programmers to create `generators` or`cli tools` for boilerplate code.
Plis uses a very simple approach to create generators all the hard work is done behind the scene from plis you only need
to use the API provided by plis and you are good to go.
## Hows does plis work ?
Plis is kind of like "npm", it has defined directories where packages (generators) are stored, generators can be 
global or local depending on how you install them. Plis reads these generators and makes it possible for you to run them.
You can simply run a generator by running:
```bash
plis {generator-name}
```
## What is a plis generator/cli tool?
A plis generator is a folder with at least 2 files:

1. A config file `config.json`
2. A run script `run.lua` (lua is one of the supported scripting languages)

Generators are stored in 2 locations:

 - Global generators :
 Global generators are stored in   `$HOME/.plis/generators`
 - Project generators:
 Project generators are stored in `plis/generators` folder in the current directory.
 
Every generator is located in a folder named `plis-{generator-name}` in one of the locations mentioned above.
 
## How do I install a generator ?

Generators are usually stored as a git repository and you can install generators by running:

 ```bash
plis get https://github.com/kujtimiihoxha/plis-generator
```

where `https://github.com/kujtimiihoxha/plis-generator` is the repository link.

*Plis uses git so git is a requirement of plis*

One other way to install generators is by providing a `plis.json` file where all the generators dependencies are listed.
```json
{
    "dependencies": [
        {
            "rep": "git@github.com:kujtimiihoxha/plis-generator",
            "branch": "master"
        }
    ]
}
```
Plis also supports a way of "versioning" by using branches so for e.x if you would like to get a specific
version of the generator you would specify the branch you want to install (assuming that the creator of the generator follows
this system of versioning), this can be done by using the `branch` property as shown in the json above or by using 
the `--branch` flag e.x:
```bash
plis get https://github.com/kujtimiihoxha/plis-generator --branch v1
```
## Give me the "Hello world" plis generator!
A simple generator has this folder structure:
```bash
├── plis-hello/
│   ├── .gitignore
│   ├── .plis-generator.json
│   ├── config.json
│   ├── run.lua
│   └── test-project/
```
- `.plis-generator.json` is used for debugging, if plis is run inside the `polis-hello` folder it knows that it is run
within a generator and sets up debugging.

*.plis-generator.json*
```json
{
    "generator_name":"hello",
    "test_dir":"test-project"
}
```
- `config.json` is used to keep all the generators configurations (name,flags,arguments, etc..)

*config.json*
```json
{
    "name":"hello",
    "description":"This is the 'hello' plis generator ",
    "script_type":"lua"
}
```
- `run.lua` is the script that is run when the generator is called.

*run.lua*
```lua
-- The main function called by plis.
function main()
    print("Hello World!")
end
```

Now the only thing you have to do is run:
```bash
plis hello
```

and the generator prints out :
```bash
Hello World!
```

- `test-project/` this folder is used when you are testing your generator, every time you run plis inside a generator
the generator that is run thinks that `test-project/` is the folder you are running plis from (the base folder),
this allows you to test your generators when you have to deal with files read,write and templates.
