# mac
 Simple CLI to replace all my random bash scripts i use on all my machines


## Getting Started

a yaml config file should be created at `~/.mac/config`

a complete example of all configuration. all configuration has "sensible" defaults

```yaml
CodeRoot: /var/www # Where code should cloned be to
QuitExclusions: #apps that should not be killed when quit is called
  - terminal
  - spotify
```