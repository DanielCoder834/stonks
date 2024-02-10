package main

import (
	"fmt"
	"github.com/google/btree"
	"strconv"
	"sync"
	"time"
)

// SETs Is one layer of a Hashmap
var SETs = map[string]string{}

// SETsMu makes sure writes don't happen at the same time of other go routines (aka cooler threads)
var SETsMu = sync.RWMutex{}

// Handlers is a map of commands to function to do each action
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	// For when connecting with non-custom redis client
	//"COMMAND":      ping,
	"SET":          set,
	"GET":          get,
	"HSET":         hset,
	"HGET":         hget,
	"DEEPGET":      deepget,
	"DEEPSETO.W":   deepsetoverwrite,
	"DEEPSETN.O.W": deepsetnotoverwrite,
	"DEEPADDLAYER": deepaddlayer,
	// Key: slices of strings, Value: strings
	"ARRSET": arrset,
	"ARRGET": arrget,
	// Data is defined to introduce a sense of time stamps
	"DataB*SET": databstarset,
	"DataB*GET": databstarget,
	// TODO: Implement a delete method
	//"ARRBGET":      arrbstarget,
	//"HGETALL": hgetall,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// ARRSETs is a hashmap with keys as string slices with the value as a string
// TODO: Find a way to make it a map for string slices instead of any interface
var ARRSETs = map[interface{}]string{}
var ARRSETsMu = sync.RWMutex{}

// First elements are the key, and the last value is the value
func arrset(args []Value) Value {
	fmt.Println(args)
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'arrset' command"}
	}

	//fmt.Println("saet works")

	//keychain := make([]string, 0)
	keychain := ""
	for argIdx := 0; argIdx < len(args)-1; argIdx++ {
		keychain += args[argIdx].bulk + "_"
	}
	value := args[len(args)-1].bulk
	ARRSETsMu.Lock()
	ARRSETs[keychain] = value
	ARRSETsMu.Unlock()

	fmt.Println("Key: ", keychain)
	fmt.Println("Value: ", value)
	return Value{typ: "string", str: "OK"}
}

func arrget(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	fmt.Println("get works")
	//keychain := make([]string, 0)
	keychain := ""
	for argIdx := 0; argIdx < len(args); argIdx++ {
		keychain += args[argIdx].bulk + "_"
	}

	ARRSETsMu.RLock()
	value, ok := ARRSETs[keychain]
	ARRSETsMu.RUnlock()

	fmt.Println("Key: ", keychain)
	fmt.Println("Value: ", value)

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// HSETs is a 2 layer hashmap
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

// Sets an item 2 layers down
func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

// Gets an item 2 layers down
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// DEEPSETs is a multiple layer hashmap, idk how deep it goes,
// TODO: Find a way to make each key unique have different values
var DEEPSETs = map[string]interface{}{}
var DEEPSETsMu = sync.RWMutex{}

func deepget(args []Value) Value {
	keychain := make([]string, 0)
	for argIdx := 0; argIdx < len(args); argIdx++ {
		keychain = append(keychain, args[argIdx].bulk)
	}
	DEEPSETsMu.RLock()
	val, ok := getdeepestvalofkeychain(DEEPSETs, keychain, 0)
	DEEPSETsMu.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

/**
 * Since this is a complex traversal, I will explain it in detail
 * The map is for recursion, it is potential next layer, e.g. map["key"] = newmap if the keys are right, then pass in newmap
 * The keychain represents a slice/list of string keys of each recursive layer such as "key" in the last example.
 * The keychain index represents the order of the slice, we consider order to be important,
 * each recurisive call the index gets incremented by one duh
 * It returns the value if found, with true or a value with typ null and false if not found
 * https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
 * https://stackoverflow.com/questions/32611829/iterate-through-a-map-of-interface-that-contains-different-levels-of-maps
 */
func getdeepestvalofkeychain(m map[string]interface{}, keychain []string, keychainidx int) (string, bool) {
	for k, v := range m {
		// Check that the key is equal to the given keychain
		if k == m[keychain[keychainidx]] {
			// Check if it is of type map, if so traverse
			if v, ok := v.(map[string]interface{}); ok {
				addOneMoreKeyChainIDx := keychainidx + 1
				getdeepestvalofkeychain(v, keychain, addOneMoreKeyChainIDx)
			} else {
				// if the type is not a map, the value has been found, yay
				// double-check the type is a value, can never be too sure
				finalstrkey := keychain[keychainidx]
				finalval := m[finalstrkey]
				if val, ok := finalval.(string); ok {
					return val, true
				}
			}
		}
	}
	return "", false
}

func deepsetoverwrite(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'deepsetoverwrite' command"}
	}
	keychain := make([]string, 0)
	valuetoset := args[len(args)].bulk
	keytoset := args[len(args)-1].bulk
	for argIdx := 0; argIdx < len(args)-2; argIdx++ {
		keychain = append(keychain, args[argIdx].bulk)
	}
	DEEPSETsMu.RLock()
	ok := setvalofkeychain(DEEPSETs, keychain, 0, valuetoset, keytoset, true)
	DEEPSETsMu.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: valuetoset}
}

func deepsetnotoverwrite(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'deepsetnotoverwrite' command"}
	}
	keychain := make([]string, 0)
	valuetoset := args[len(args)].bulk
	keytoset := args[len(args)-1].bulk
	for argIdx := 0; argIdx < len(args)-2; argIdx++ {
		keychain = append(keychain, args[argIdx].bulk)
	}
	DEEPSETsMu.RLock()
	ok := setvalofkeychain(DEEPSETs, keychain, 0, valuetoset, keytoset, false)
	DEEPSETsMu.Unlock()
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: valuetoset}
}

/**
 * Also a bit of a complex function don't worry, same as the rest, recursively traverse the hashnmaps until the keychain ends
 * once it ends, set the desired key to the desired value.
 * the overwrite bool gives the option to look for if the key exists at the end of the hashmap levels
 * if the overwrite is set true, we don't care if the value exists or not
 */
func setvalofkeychain(m map[string]interface{}, keychain []string, keychainidx int, valtoset string, keytoset string, overwrite bool) bool {
	for k, v := range m {
		// Check if the end of the recursive step
		if keychainidx == len(keychain)+1 {
			// If overwrite we don't care if there is a value or not
			if overwrite {
				m[keytoset] = valtoset
				return true
			} else if !overwrite && !keyinmap(m, keytoset) {
				m[keytoset] = valtoset
				return true
			} else {
				return false
			}
		}
		// Check that the key (in the loop) is equal to the given keychain
		if k == m[keychain[keychainidx]] {
			// Check if it is of type map, if so traverse
			if v, ok := v.(map[string]interface{}); ok {
				addOneMoreKeyChainIDx := keychainidx + 1
				setvalofkeychain(v, keychain, addOneMoreKeyChainIDx, valtoset, keytoset, overwrite)
			}
		}
	}
	return false
}

func keyinmap(m map[string]interface{}, keytoset string) bool {
	for key := range m {
		if key == keytoset {
			return true
		}
	}
	return false
}

var LAYERsMu = sync.RWMutex{}

// "Multiple-level hashmaps are like onions. Onions have layers. Multiple-level hashmaps have layers." - Shrek
func deepaddlayer(args []Value) Value {
	keychain := make([]string, 0)
	keytoset := args[len(args)].bulk
	for argIdx := 0; argIdx < len(args)-1; argIdx++ {
		keychain = append(keychain, args[argIdx].bulk)
	}
	LAYERsMu.Lock()
	ok := recursiveLayerTraversalAndAdding(DEEPSETs, keychain, 0, keytoset)
	LAYERsMu.Unlock()
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: keytoset}
}

/**
 * Same as the rest, recursively traverse the multiple layer hasmap until the keychain ends
 * The m represents the map, the keychain represents the list of keys to traverse the map, the keychain index represents the index of the keychain, and
 * the key to set represents the desire key of the new hashmap
 * Once the len of the keychain is equal to the keychain index, the keychain is done, which means the new layer gets set and the returns true
 * In the case false is returned, the keychain is invalid or I screwed up
 */
func recursiveLayerTraversalAndAdding(m map[string]interface{}, keychain []string, keychainidx int, keytoset string) bool {
	for k, v := range m {
		if k == m[keychain[keychainidx]] {
			// Check if it is of type map, if so traverse
			if v, ok := v.(map[string]interface{}); ok {
				if keychainidx == len(keychain) {
					v[keytoset] = map[string]interface{}{}
					return true
				} else {
					addOneMoreKeyChainIDx := keychainidx + 1
					recursiveLayerTraversalAndAdding(v, keychain, addOneMoreKeyChainIDx, keytoset)
				}
			}
		}
	}
	return false
}

type DataAges string

const (
	Fresh DataAges = "Fresh"
	Valid DataAges = "Valid"
	Stale DataAges = "Stale"
)

type Data struct {
	Time time.Time
	Value
	DataAges DataAges
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func (a Data) Less(b btree.Item) bool {
	castItemTime := b.(*Data)
	return inTimeSpan(a.Time, time.Now(), castItemTime.Time)
}

var degrees = 6
var BTREESETs = btree.New(degrees)
var BTREEMu = sync.RWMutex{}

func arrbstarget(values []Value) Value {
	//BTREESETs.ReplaceOrInsert(values[0])
	fmt.Println("get8: ", BTREESETs.Get(btree.Int(8)))
	//item := btree.Item(8);
	return Value{typ: "null"}
}

func databstarset(values []Value) Value {
	if len(values) != 8 {
		return Value{typ: "null"}
	}
	year, yearErr := strconv.Atoi(values[0].bulk)
	if yearErr != nil {
		return Value{typ: "null"}
	}
	monthInt, montherr := strconv.Atoi(values[1].bulk)
	if montherr != nil {
		return Value{typ: "null"}
	}
	month := time.Month(monthInt)

	value := values[8]
	//value := values[8].bulk
	//value := values[8].bulk
	//value := values[8].bulk
	//value := values[8].bulk
	//value := values[8].bulk
	data := Data{
		DataAges: Fresh,
		Value:    value,
		Time:     time.Date(year, month, 0, 0, 0, 0, 0, time.UTC),
	}
	return value
}

func databstarget(values []Value) Value {
}

//func numoflayers(m map[string]interface{}) map[string]interface{} {
//	for k, v := range m {
//		_ = k
//		if v, ok := v.(map[string]interface{}); ok {
//			getdeepestmap(v)
//		}
//	}
//	return m
//}

//func bstarget(args []Value) Value {
//keychain := []string
//for idxArgs := 0; idxArgs < len(args) - 1; idxArgs++ {
//	append(keychain, args[idxArgs].bulk)
//}

//argsItem := args.(btree.Item)
//retrievedItem := BTREESETs.Get(argsItem)
//convertToValue := retrievedItem.([]Value)
//return convertToValue
//}
