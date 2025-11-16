#!/bin/bash
# Fix ZSH Terminal - Run this script to restore terminal functionality

echo "🔧 Fixing ZSH configuration..."

# Backup current config
cp ~/.zshrc ~/.zshrc.backup.$(date +%Y%m%d_%H%M%S)
echo "✅ Backed up current .zshrc"

# Create minimal working config
cat > ~/.zshrc << 'EOF'
#=================================================
# Minimal Zsh Configuration - Emergency Fix
#=================================================

# Basic PATH
export PATH="$HOME/bin:$HOME/.local/bin:$HOME/.cargo/bin:/usr/local/bin:/usr/bin:/bin"

# Basic aliases
alias ll='ls -lah'
alias la='ls -A'
alias l='ls -CF'

# History
HISTSIZE=1000
SAVEHIST=1000
HISTFILE=$HOME/.zsh_history
setopt SHARE_HISTORY HIST_IGNORE_DUPS

# Basic completion
autoload -Uz compinit
compinit

# Simple prompt
PS1='%n@%m:%~$ '

# NVM (if exists)
if [ -d "$HOME/.nvm" ]; then
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
fi

echo "✅ ZSH loaded successfully (minimal config)"
EOF

echo "✅ Created minimal .zshrc"
echo ""
echo "📝 To restore full config later:"
echo "   mv ~/.zshrc.backup.YYYYMMDD_HHMMSS ~/.zshrc"
echo ""
echo "🚀 Restart your terminal or run: source ~/.zshrc"