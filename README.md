# Token Ring Emulator (Go)

Эмулятор протокола Token Ring, реализованный на языке Go с использованием
исключительно стандартной библиотеки.

## Описание

Программа моделирует кольцевую сеть из N узлов.  
Каждый узел работает в отдельной goroutine и связан каналами с соседними узлами.

Сообщения (Token) передаются по кольцу от узла к узлу, пока:
- не достигнут адресата (по sha3-хэшу идентификатора узла), либо
- не истечёт время жизни сообщения (TTL).

Первое сообщение отправляется главным потоком узлу №1.
Каждое последующее сообщение генерируется узлом, получившим сообщение.

## Формат Token

Каждое сообщение содержит:
- данные (строка),
- sha3-хэш идентификатора узла-получателя,
- счётчик времени жизни (TTL).

## Требования

- Go версии 1.25 или выше
- Используется только стандартная библиотека Go

## Запуск программы

```bash
go run tokenring.go <number_of_nodes>

## Пример вывода

[Main] sending initial token → 4
[Node 2] received token: data="initial message" ttl=4
[Node 3] forwarding token
[Node 4] token accepted: "initial message"

[Main] sending initial token → 2
[Node 2] received token: data="initial message" ttl=2
[Node 2] token accepted: "initial message"
[Node 2] generated new token → 1
[Node 1] received token: data="msg from 2 to 1" ttl=2
[Node 1] token accepted: "msg from 2 to 1"
[Node 1] generated new token → 1
[Node 2] received token: data="msg from 1 to 1" ttl=1