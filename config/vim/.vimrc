set number
set relativenumber
set mouse=a
set cursorline
set noswapfile
set nobackup
set nocompatible

set scrolloff=8
set smarttab
set expandtab
set tabstop=2 
set shiftwidth=2
set softtabstop=2
set autoindent

set t_Co=256
syntax on
set mousehide
set termencoding=utf-8
set wrap
set linebreak
set visualbel t_vb=

set encoding=utf-8
set fileencodings=utf-8,cp1251

set clipboard=unnamedplus
set ruler
set hidden
nnoremap <S-n> :bnext<CR>
nnoremap <S-p> :bprevious<CR>
inoremap jj <Esc>
nnoremap ; :

set wildmenu
filetype plugin indent on

nnoremap <leader>e :Explore<CR>

let mapleader = " "

let g:netrw_banner = 0
let g:netrw_listyle = 3
let g:netrw_browse_split = 3
