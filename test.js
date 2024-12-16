const createPosts = async (title, body) => {
    const url = "http://10.1.9.2:8080/posts/store"; // Replace with your actual server URL

    for (let i = 1; i <= 10; i++) {
        const formData = new URLSearchParams();
        formData.append("title", `${title} ${i}`);
        formData.append("body", `${body} ${i}`);

        try {
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Cookie": "token=8f5209d6-b7d5-11ef-a6ec-047c1698dc8d"
                },
                body: formData
            });

            if (!response.ok) {
                console.error(`Failed to create post ${i}:`, response.statusText);
            } else {
                console.log(`Post ${i} created successfully!`);
            }
        } catch (error) {
            console.error(`Error creating post ${i}:`, error);
        }
    }
};

// Example usage
const title = "Post Title";
const body = "This is the body of the post.";
createPosts(title, body);