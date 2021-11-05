package chunkyUploads

// type Chunk struct {
// 	UUID UUID
// 	Hash Hash
// }

type ChunkAttachment struct {
	// File       UUID
	Chunk UUID
	Start int64
	End   int64
}
