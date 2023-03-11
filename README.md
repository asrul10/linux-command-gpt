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

### Example Usage

```bash
> lcg I want to extract file linux-command-gpt.tar.gz
Completed in 0.92 seconds
┌────────────────────────────────────┐
│ tar -xvzf linux-command-gpt.tar.gz │
└────────────────────────────────────┘
Are you sure you want to execute the command? (Y/n):
```

### Options
```bash
> lcg [options]

--help         output usage information
--version      output the version number
--update-key   update the API key
--delete-key   delete the API key
```
