package tests

//Introduction to testing.  Note that testing is built into go and we will be using
//it extensively in this class. Below is a starter for your testing code.  In
//addition to what is built into go, we will be using a few third party packages
//that improve the testing experience.  The first is testify.  This package brings
//asserts to the table, that is much better than directly interacting with the
//testing.T object.  Second is gofakeit.  This package provides a significant number
//of helper functions to generate random data to make testing easier.

import (
	"fmt"
	"os"
	"testing"

	"drexel.edu/todo/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/stretchr/testify/assert"
)

// Note the default file path is relative to the test package location.  The
// project has a /tests path where you are at and a /data path where the
// database file sits.  So to get there we need to back up a directory and
// then go into the /data directory.  Thus this is why we are setting the
// default file name to "../data/todo.json"
const (
	DEFAULT_DB_FILE_NAME = "../data/todo.json"
)

var (
	DB *db.ToDo
)

// note init() is a helpful function in golang.  If it exists in a package
// such as we are doing here with the testing package, it will be called
// exactly once.  This is a great place to do setup work for your tests.
func init() {
	//Below we are setting up the gloabal DB variable that we can use in
	//all of our testing functions to make life easier
	testdb, err := db.New(DEFAULT_DB_FILE_NAME)
	if err != nil {
		fmt.Print("ERROR CREATING DB:", err)
		os.Exit(1)
	}

	DB = testdb //setup the global DB variable to support test cases

	//Now lets start with a fresh DB with the sample test data
	testdb.RestoreDB()
}

// Sample Test, will always pass, comparing the second parameter to true, which
// is hard coded as true
func TestTrue(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestAddHardCodedItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     999,
		Title:  "This is a test case item",
		IsDone: false,
	}
	t.Log("Testing Adding a Hard Coded Item: ", item)

	//TODO: finish this test, add an item to the database and then
	//check that it was added correctly by looking it back up
	//use assert.NoError() to ensure errors are not returned.
	//explore other useful asserts in the testify package, see
	//https://github.com/stretchr/testify.  Specifically look
	//at things like assert.Equal() and assert.Condition()

	//I will get you started, uncomment the lines below to add to the DB
	//and ensure no errors:
	//---------------------------------------------------------------
	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to DB")

	//TODO: Now finish the test case by looking up the item in the DB
	//and making sure it matches the item that you put in the DB above
	dbItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error getting item from DB")

	assert.Equal(t, item, dbItem, "Items don't match")
}

func TestAddRandomStructItem(t *testing.T) {
	//You can also use the Stuct() fake function to create a random struct
	//Not going to do anyting
	item := db.ToDoItem{}
	err := fake.Struct(&item)
	t.Log("Testing Adding a Randomly Generated Struct: ", item)

	assert.NoError(t, err, "Created fake item OK")

	//TODO: Complete the test
	err = DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to DB")

	dbItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error getting item from DB")

	assert.Equal(t, item, dbItem, "Items don't match")
}

func TestAddRandomItem(t *testing.T) {
	//Lets use the fake helper to create random data for the item
	item := db.ToDoItem{
		Id:     fake.Number(100, 110),
		Title:  fake.JobTitle(),
		IsDone: fake.Bool(),
	}

	t.Log("Testing Adding an Item with Random Fields: ", item)
	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item to DB")

	dbItem, err := DB.GetItem(item.Id)
	assert.NoError(t, err, "Error getting item from DB")

	assert.Equal(t, item, dbItem, "Items don't match")
}

//TODO: Create additional tests to showcase the correct operation of your program
//for example getting an item, getting all items, updating items, and so on. Be
//creative here.

func TestAddItemWithExistingId(t *testing.T) {
	item := db.ToDoItem{
		Id:     2,
		Title:  "This is a test case item",
		IsDone: false,
	}

	err := DB.AddItem(item)

	assert.EqualError(t, err, "Item id already exists in db")
}

func TestGetItem(t *testing.T) {
	expectedItem := db.ToDoItem{
		Id:     1,
		Title:  "Learn Go / GoLang",
		IsDone: false,
	}

	actualItem, err := DB.GetItem(1)
	assert.NoError(t, err, "Error getting item.")

	assert.Equal(t, expectedItem, actualItem)
}

func TestGetItemNoItem(t *testing.T) {
	_, err := DB.GetItem(10)

	assert.EqualError(t, err, "Id not found in db.")
}

func TestDeleteItemNoItem(t *testing.T) {
	err := DB.DeleteItem(10)

	assert.EqualError(t, err, "Item id does not exist in db")
}

func TestDeleteItem(t *testing.T) {
	item := db.ToDoItem{
		Id:     998,
		Title:  "This is a test case item",
		IsDone: false,
	}

	err := DB.AddItem(item)
	assert.NoError(t, err, "Error adding item.")

	err = DB.DeleteItem(998)
	assert.NoError(t, err, "Error deleting item")

	_, err = DB.GetItem(998)
	assert.EqualError(t, err, "Id not found in db.")
}

func TestUpdateItemWhenMissing(t *testing.T) {
	updatedItem := db.ToDoItem{
		Id:     10,
		Title:  "Learn Go / GoLang",
		IsDone: true,
	}
	err := DB.UpdateItem(updatedItem)
	assert.EqualError(t, err, "Item id does not exist in db")
}

func TestUpdateItem(t *testing.T) {
	updatedItem := db.ToDoItem{
		Id:     1,
		Title:  "Learn Go / GoLang",
		IsDone: true,
	}
	err := DB.UpdateItem(updatedItem)
	assert.NoError(t, err, "Error updating item.")

	actualItem, err := DB.GetItem(1)
	assert.NoError(t, err, "Error getting item.")
	assert.Equal(t, updatedItem, actualItem, "Items don't match after update.")
}

func TestGetAllItems(t *testing.T) {
	items, err := DB.GetAllItems()
	assert.NoError(t, err, "Error occurred getting all items")

	assert.GreaterOrEqual(t, len(items), 4, "There should be at least 4 items.")

	for id := 1; id <= 4; id++ {
		item, err := DB.GetItem(id)
		assert.NoErrorf(t, err, "Error getting item %d", id)
		assert.Containsf(t, items, item, "List does not contain item %d", id)
	}
}

func TestChangeItemDoneStatusWhenItemMissing(t *testing.T) {
	err := DB.ChangeItemDoneStatus(10, true)
	assert.EqualError(t, err, "Id not found in db.")
}

func TestChangeItemDoneStatus(t *testing.T) {
	err := DB.ChangeItemDoneStatus(1, true)
	assert.NoError(t, err, "No error updating item.")

	item, err := DB.GetItem(1)
	assert.NoError(t, err, "Error getting item.")

	assert.Equal(t, true, item.IsDone, "Item status not updated.")

}
