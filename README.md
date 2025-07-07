# gh-commitmsg

Идея приложения:
- получить staged changes в git репозитории в текущей директории
- скормить полученные данные LLM GitHub Models для генерации conventional commit message
- вывести полученный commit message на экран
- если будет хорошо получаться, можно в качестве примера для модели давать несколько предыдущих commit messages