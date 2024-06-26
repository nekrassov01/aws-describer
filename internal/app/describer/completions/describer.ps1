# This script inspired by https://github.com/urfave/cli
# NOTE: Complex completions such as flag combination checks are not supported

Register-ArgumentCompleter -Native -CommandName "describer" -ScriptBlock {
  param($commandName, $wordToComplete, $cursorPosition)
  (Invoke-Expression "$wordToComplete --generate-bash-completion").ForEach{
    [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
  }
}
