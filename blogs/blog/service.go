package blog

// BlogService handles the blog-related business logic.
type BlogService struct{}

// Create a new blog post using the Blog model.
func (s *BlogService) CreateBlog(blog Blog) error {
	// You can add additional business logic here if needed
	return blog.CreateBlog()
}

// Retrieve all blog posts.
func (s *BlogService) GetAllBlogs() ([]Blog, error) {
	// Additional logic can be added here, e.g., filtering, ordering, pagination
	return GetAllBlogs()
}

// Retrieve a single blog post by its ID.
func (s *BlogService) GetBlogByID(id int) (*Blog, error) {
	return GetBlogByID(id)
}

// Update an existing blog post.
func (s *BlogService) UpdateBlog(blog Blog) error {
	// Additional logic for updates (e.g., checking permissions, validating input) can be added here
	return blog.UpdateBlog()
}

// Delete a blog post by its ID.
func (s *BlogService) DeleteBlog(id int) error {
	// Call the model's DeleteBlog function to remove the blog post
	return DeleteBlogFromModel(id)
}
