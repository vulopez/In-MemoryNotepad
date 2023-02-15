package main

import (
	"bufio"
	. "fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	startMsg         = "Enter the maximum number of notes: "
	inpMsg           = "Enter a command and data: "
	addMsg           = "[OK] The note was successfully created"
	delMsg           = "[OK] All notes were successfully deleted"
	updateMsg        = "[OK] The note at position %d was successfully updated"
	deleteMsg        = "[OK] The note at position %d was successfully deleted"
	fullMsg, listFmt = "[Error] Notepad is full", "[Info] %d: %s"
	emptyMsg         = "[Info] Notepad is empty"
	emptyPos         = "[Error] Missing position argument"
	exitMsg          = "[Info] Bye!"
	errCmd           = "[Error] Unknown command"
	errVal           = "[Error] Missing note argument"
	errPos           = "[Error] Invalid position: %s"
	errUpd           = "[Error] There is nothing to update"
	errDel           = "[Error] There is nothing to delete"
	errOutBounds     = "[Error] Position %d is out of the boundaries [1, %d]"
)

var (
	arrLimit int // absolute cap(notepad)
	scanner  = bufio.NewScanner(os.Stdin)
)

func initNotepad() []string {
	sl := make([]string, 0, arrLimit)
	notes, _ := clear(sl)
	return notes
}

func splitMsg(data, content string) (string, string) {
	entry := regexp.MustCompile(`\s+`).Split(content, 2)
	entry = append(entry, data)
	return entry[0], entry[1]
}

func add(notes []string, content string) ([]string, string) {
	if content == "" {
		return notes, errVal
	}
	if len(notes) == arrLimit {
		return notes, fullMsg
	}
	notes = append(notes, content)
	return notes, addMsg
}

func list(notes []string) string {
	var listed = initNotepad()
	if notes != nil {
		for i, note := range notes {
			listed = append(listed, Sprintf(listFmt, i+1, note))
		}
		return strings.Join(listed, "\n")
	}
	return emptyMsg
}

func clear(notes []string) ([]string, string) {
	notes = nil
	return notes, delMsg
}

func updateNote(notes []string, prefix string, content, data string) ([]string, string) {
	pos, err := strconv.Atoi(prefix)
	if prefix == "" { // empty position
		return notes, emptyPos
	}
	if len(strings.Fields(data)) == 1 { // empty message
		return notes, errVal
	}
	if err != nil { // cannot convert to Int
		return notes, Sprintf(errPos, prefix)
	}
	if arrLimit < pos { // out boundy
		return notes, Sprintf(errOutBounds, pos, arrLimit) // errDel
	}
	if len(notes) < pos { // Nothing to delete
		return notes, Sprintf(errUpd)
	}
	notes[pos-1] = content
	return notes, Sprintf(updateMsg, pos)
}

func deleteNote(notes []string, prefix string) ([]string, string) {
	tempNotepad := initNotepad()
	pos, err := strconv.Atoi(prefix)
	if prefix == "" { //
		return notes, Sprintf(emptyPos)
	}
	if err != nil { // cannot convert to Int
		return notes, Sprintf(errPos, prefix)
	}
	if len(notes) < pos || notes[pos] == "" {
		return notes, Sprintf(errDel) // errDel
	}
	if arrLimit < pos {
		return notes, Sprintf(errOutBounds, pos, arrLimit) // errDel
	}

	notes[pos-1] = ""
	for _, note := range notes {
		if note != "" {
			tempNotepad = append(tempNotepad, note)
		}
	}
	return tempNotepad, Sprintf(deleteMsg, pos)
}

func main() {
	var command string
	Printf("%s", startMsg)
	// get notepad size
	_, _ = Scanf("%d\n", &arrLimit)

	notepad := initNotepad()
	// TODO: improve this break line
	Println()
mainLoop:
	for {
		var data, response string
		Printf("%s", inpMsg)
		scanner.Scan()

		command, data = splitMsg(data, scanner.Text())

		switch command {
		case "create":
			notepad, response = add(notepad, data)
		case "list":
			response = list(notepad)
		case "clear":
			notepad, response = clear(notepad)
		case "update":
			notePos, msg := splitMsg(data, data)
			notepad, response = updateNote(notepad, notePos, msg, data)
		case "delete":
			notepad, response = deleteNote(notepad, data)
		case "exit":
			break mainLoop
		default:
			// Print(errCmd)
			response = errCmd
		}
		Println(response)
		// TODO: delete this break line
		Println()
	}
	Printf(exitMsg)
}
