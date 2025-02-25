-- opts
vim.g.loaded_perl_provider = 0
vim.g.loaded_ruby_provider = 0
vim.g.loaded_python_provider = 0

local opts = vim.opt

opts.number = true
opts.relativenumber = true

opts.swapfile = false
opts.encoding = "utf-8"
opts.clipboard = "unnamedplus"

opts.scrolloff = 8

opts.tabstop = 2
opts.softtabstop = 2
opts.shiftwidth = 2
opts.expandtab = true
opts.autoindent = true
opts.cursorline = true

vim.g.netrw_banner = 0
vim.g.netrw_liststyle = 3
vim.g.netrw_browse_split = 3

vim.cmd("colorscheme desert")

opts.spell = false
vim.o.showmode = false
vim.o.laststatus = 2
vim.o.statusline = "%!v:lua.statusline()"

-- autocmd
function _G.statusline()
	local bufname = vim.fn.expand("%:t")
	local bufnr = vim.fn.bufnr("%")
	local buf_count = vim.fn.len(vim.fn.getbufinfo("%"))
	return string.format(" %s | Buffers: %d ", bufname, buf_count)
end

-- keymap
local keyset = vim.keymap.set
local opts = { noremap = true, silent = true }

vim.g.mapleader = " "
vim.g.maplocalleader = ","

keyset("n", ";", ":")
keyset("i", "jj", "<Esc>", { desc = "Exit from insert" })
keyset("n", "<leader>sh", ":split<CR>", opts, { desc = "Horizontal split" })
keyset("n", "<leader>ss", ":vsplit<CR>", opts, { desc = "Vertical split" })
keyset("n", "<Esc>", ":nohlsearch<CR>", opts)

keyset("n", "<leader>e", ":Ex<CR>", opts, { desc = "Explorer" })
keyset("n", "<C-a>", "gg<S-v>G", { desc = "Select all text" })

keyset("n", "<C-h>", "<C-W>h", opts)
keyset("n", "<C-j>", "<C-W>j", opts)
keyset("n", "<C-k>", "<C-W>k", opts)
keyset("n", "<C-l>", "<C-W>l", opts)

keyset("n", "<S-h>", "<Cmd>bprev<CR>", opts)
keyset("n", "<S-l>", "<Cmd>bnext<CR>", opts)
keyset("n", "<leader>bc", "<Cmd>bdelete<CR>", opts)
vim.api.nvim_set_keymap("n", "<C-b>", ":enew<CR>", { noremap = true, silent = true })
