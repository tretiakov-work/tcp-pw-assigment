package quote_generator

import (
	"math/rand"
)

type StaticMemory struct {
	quotes [][]byte
}

func NewStaticMemory() *StaticMemory {
	return &StaticMemory{
		quotes: [][]byte{
			[]byte("Be the change you wish to see in the world. — Mahatma Gandhi"),
			[]byte("In the middle of every difficulty lies opportunity. — Albert Einstein"),
			[]byte("The only limit to our realization of tomorrow is our doubts of today. — Franklin D. Roosevelt"),
			[]byte("You must be the change you wish to see in the world. — Mahatma Gandhi"),
			[]byte("The journey of a thousand miles begins with one step. — Lao Tzu"),
			[]byte("What lies behind us and what lies before us are tiny matters compared to what lies within us. — Ralph Waldo Emerson"),
			[]byte("The only way to do great work is to love what you do. — Steve Jobs"),
			[]byte("Life is what happens when you're busy making other plans. — John Lennon"),
			[]byte("The best way to predict the future is to invent it. — Alan Kay"),
			[]byte("Your time is limited, don't waste it living someone else's life. — Steve Jobs"),
			[]byte("The only impossible journey is the one you never begin. — Tony Robbins"),
			[]byte("Success is not the key to happiness. Happiness is the key to success. — Albert Schweitzer"),
			[]byte("Believe you can and you're halfway there. — Theodore Roosevelt"),
			[]byte("Life is either a daring adventure or nothing at all. — Helen Keller"),
			[]byte("The best revenge is massive success. — Frank Sinatra"),
			[]byte("Act as if what you do makes a difference. It does. — William James"),
			[]byte("In three words I can sum up everything I've learned about life: it goes on. — Robert Frost"),
			[]byte("You only live once, but if you do it right, once is enough. — Mae West"),
			[]byte("The purpose of our lives is to be happy. — Dalai Lama"),
			[]byte("You miss 100% of the shots you don't take. — Wayne Gretzky"),
			[]byte("The only limit to our realization of tomorrow is our doubts of today. — Franklin D. Roosevelt"),
			[]byte("Don't watch the clock; do what it does. Keep going. — Sam Levenson"),
			[]byte("Life isn't about finding yourself. Life is about creating yourself. — George Bernard Shaw"),
			[]byte("To be yourself in a world that is constantly trying to make you something else is the greatest accomplishment. — Ralph Waldo Emerson"),
			[]byte("What lies behind us and what lies before us are tiny matters compared to what lies within us. — Ralph Waldo Emerson"),
			[]byte("The only way to do great work is to love what you do. — Steve Jobs"),
			[]byte("The best way to find yourself is to lose yourself in the service of others. — Mahatma Gandhi"),
			[]byte("To live is the rarest thing in the world. Most people exist, that is all. — Oscar Wilde"),
			[]byte("The mind is everything. What you think you become. — Buddha"),
			[]byte("Do not wait to strike till the iron is hot, but make it hot by striking. — William Butler Yeats"),
		},
	}
}

func (q *StaticMemory) GetQuote() ([]byte, error) {
	rnd := rand.Intn(len(q.quotes) - 1)
	return q.quotes[rnd], nil
}
