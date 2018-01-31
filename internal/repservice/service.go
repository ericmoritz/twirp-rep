package repservice

import (
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	pb "github.com/ericmoritz/twirp-rep/rpc/rep"
	usersRPC "github.com/ericmoritz/twirp-users/rpc/users"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)


type repService struct {
	DB *leveldb.DB
	UsersClient usersRPC.Users
}

// Creates a new Rep service
func New(dbPath string, usersClient usersRPC.Users) (pb.Rep, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &repService{
		DB: db,
		UsersClient: usersClient,
	}, nil
}

func (rs *repService) Vote(c context.Context, req *pb.VoteReq) (*pb.VoteResp, error) {
	// Attempt to get the current user
	currentUserResp, err := rs.UsersClient.CurrentUser(
		context.Background(),
		&usersRPC.CurrentUserReq{Session: req.Session},
	)
	if err != nil {
		return nil, err
	}
	// Create a vote message
	vote := &pb.VoteRecord{
		UUID: uuid.NewV4().String(),
		Source: currentUserResp.User.Username,
		Target: req.Target,
	}

	// Store the vote
	err = putVote(rs.DB, vote)
	return &pb.VoteResp{}, err
}

func (rs *repService) Votes(c context.Context, req *pb.VotesReq) (*pb.VotesResp, error) {
	votes, err := countVotes(rs.DB, req.Username)
	if err != nil {
		return nil, err
	}
	return &pb.VotesResp{Count: votes}, nil
}

///////////////////////////////////////////////////////////////////////////////
// Internal
///////////////////////////////////////////////////////////////////////////////

func putVote(db *leveldb.DB, vote *pb.VoteRecord) error {
	bytes, err := proto.Marshal(vote)
	if err != nil {
		return err
	}

	return db.Put(voteKey(vote), bytes, nil)
}

func voteKey(vote *pb.VoteRecord) []byte {
	return []byte("votes/" + vote.Target + "/" + vote.Source + "/" + vote.UUID)
}

func votesKey(username string) []byte {
	return []byte("votes/" + username + "/")
}

func countVotes(db *leveldb.DB, username string) (uint64, error) {
	var count uint64
	iter := db.NewIterator(util.BytesPrefix(votesKey(username)), nil)

	for iter.Next() {
		count = count + 1
	}
	iter.Release()
	return count, iter.Error()
}
