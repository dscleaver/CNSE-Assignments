# Voting API Proposal
## Dave Cleaver

### API Structure Changes

I chose to follow a Hypermedia Application Language (HAL) approach to building these APIs as it offers a clean way to present the elements of the APIs.

#### Utility Structures

I would add the following utility structures for use across all the APIs. The Link structure would need to implement TextMarshaler to ensure that the Href properly gets properly serialized as a string. 

```go
type link struct {
  Href url.URL
}

type selfLink struct {
  Self link
}
```

#### The `Voter` API

The vote history for the voters would be modified to include their own self link and the link to the poll that the vote was for. See below:

```go
type voterPollLinks struct {
  selfLink
  Poll link
}

type voterPoll struct {
  Links   voterPollLinks `json:"_links"`
  PollID    uint
  VoteDate  time.Time
}
```

The voter item itself is changed to add the self link and embed the vote history. Also a link to the voter history is made available to provide the user with guidance on the proper url to use when adding votes. The same Voter structure can be re-used to fulfill requests for that resource. See below:

```go
type voterEmbedded struct {
  VoteHistory []voterPoll
}


type voterLinks struct {
  selfLink
  VoterHistory link
}

type Voter struct {
  Links       voterLinks    `json:"_links"`
  Embedded    voterEmbedded `json:"_embedded"`
  VoterID     uint
  FirstName   string
  LastName    string
}
```

Finally I would add an explicit structure for returning the list of votes as an embedded set of HAL Resources from `GET /voters`. See below:

```go
type votersEmbedded struct {
  Voters []Voter
}

type Voters struct {
  Links       selfLink       `json:"_links"`
  Embedded    votersEmbedded `json:"_embedded"`
}
```

#### The `Poll` API

Similar to the changes made on voters, the pollOptions become HAL resource with self links and are then embedded into the Poll Resource which also gains the self link and a direct link to the poll options. See below:

```go
type pollOption struct {
  Links           selfLink `json:"_links"`
  PollOptionID    uint
  PollOptionText  string
}

type pollEmbedded struct {
  PollOptions []pollOption
}

type pollLinks struct {
  selfLink
  PollOptions link
}

type Poll struct {
  Links        pollLinks    `json:"_links"`
  Embedded     pollEmbedded `json:"_embedded"`
  PollID       uint
  PollTitle    string
  PollQuestion string
}
```

Also much like the voter API, I addded an additional resource to return all the Polls as embedded resources from `GET /polls` See below:

```go
type pollsEmbedded {
  Polls []Poll
}

type Polls struct {
  Links    selfLink      `json:"_links"`
  Embedded pollsEmbedded `json:"_embedded"`
}
```

#### The `Votes` API

Each Vote is modified to add links to the Voter, Poll, and PollOption. See below:

```go
type voteLinks {
  selfLink
  Voter link
  Poll link
  VoteValue link
}

type Vote struct {
  Links     voteLinks `json:"_links"`
  VoteID    uint
  VoterID   uint
  PollID    uint
  VoteValue uint
}
```

Once again a resource was added to be return from `GET /votes`. See below:

```go
type votesEmbedded struct {
  Votes []Vote
}

type Votes struct {
  Links    selfLink      `json:"_links"`
  Embedded votesEmbedded `json:"_embedded"`
}
```

#### The `Home` page

For the convenience of the user and to appropriately comply with recommended HATEOAS approaches. I would also recommend providing a static resource available at `GET /` See below:

```json
{
  "_links": {
    "self": {
      "href": "/",
    },
    "voters": {
      "href": "/voters"
    },
    "polls": {
      "href": "/polls"
    },
    "votes": {
      "href": "/votes"
    }
  }
}
```

#### Notes

I chose in each of the resources to keep the various id properties in the HAL resources. One of my references was arguing against this approach, but for the purposes of this set of APIs with relatively simple links, it seemed like the best approach to make `POST` and `PUT` requests straightforward for the user.

### Example



### References

* [HAL Draft Specification](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-11)
* [Working with HAL in Put](https://evertpot.com/working-with-hal-in-put/)