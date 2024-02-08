## Linux Command GPT (lcg)
Get Linux commands in natural language with the power of ChatGPT.

### Installation
Build from source
```bash
> git clone --depth 1 https://github.com/asrul10/linux-command-gpt.git ~/.linux-command-gpt
> cd ~/.linux-command-gpt
> go build -o lcg
# Add to your environment $PATH
> ln -s ~/.linux-command-gpt/lcg ~/.local/bin
```

Or you can [download lcg executable file](https://github.com/asrul10/linux-command-gpt/releases)

### Example Usage

```bash
> lcg I want to extract linux-command-gpt.tar.gz file
Completed in 0.92 seconds

tar -xvzf linux-command-gpt.tar.gz 

Do you want to (c)opy, (r)egenerate, or take (N)o action on the command? (c/r/N):
```

To use the "copy to clipboard" feature, you need to install either the `xclip` or `xsel` package.

### Options
```bash
> lcg [options]

--help        -h  output usage information
--version     -v  output the version number
--file        -f  read command from file
--update-key  -u  update the API key
--delete-key  -d  delete the API key
```
