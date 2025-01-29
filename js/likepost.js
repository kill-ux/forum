const likesPost = document.querySelectorAll(".likesPost");

likesPost.forEach((post) => {
  post.addEventListener("submit", (e) => e.preventDefault());
  const btnLike = post.querySelector(".btn-like-post");
  const btnDisLike = post.querySelector(".btn-dislike-post");

  btnLike.addEventListener("click", async () => {
    const data = new FormData(post);
    data.append("like", "1");
    try {
      const response1 = await fetch("http://localhost:8080/posts/likes", {
        method: "POST",
        body: data,
      });
      if (response1.ok) {
        if (response1.redirected) {
          location.href = response1.url;
        } else {
          const result = await response1.json();
          const svgLike = btnLike.querySelector("svg");
          const svgDisLike = btnDisLike.querySelector("svg");
          const spanLike = btnLike.querySelector("span");
          const spanDisLike = btnDisLike.querySelector("span");
          spanLike.textContent = result.Likes;
          spanDisLike.textContent = result.Dislikes;
          if (result.Did) {
            if (result.Like) {
              svgLike.classList.add("like-blue");
              svgDisLike.classList.remove("like-red");
            }
          } else {
            svgLike.classList.remove("like-blue");
            svgDisLike.classList.remove("like-red");
          }
        }
      } else {
        throw new Error("not ok!");
      }
    } catch (error) {
      console.log(error);
    }
  });

  btnDisLike.addEventListener("click", async () => {
    const data = new FormData(post);
    data.append("like", "0");
    try {
      const response1 = await fetch("http://localhost:8080/posts/likes", {
        method: "POST",
        body: data,
      });
      if (response1.ok) {
        if (response1.redirected) {
          location.href = response1.url;
        } else {
          const result = await response1.json();
          const svgLike = btnLike.querySelector("svg");
          const svgDisLike = btnDisLike.querySelector("svg");
          const spanLike = btnLike.querySelector("span");
          const spanDisLike = btnDisLike.querySelector("span");
          spanLike.textContent = result.Likes;
          spanDisLike.textContent = result.Dislikes;
          if (result.Did) {
            if (!result.Like) {
              svgLike.classList.remove("like-blue");
              svgDisLike.classList.add("like-red");
            }
          } else {
            svgLike.classList.remove("like-blue");
            svgDisLike.classList.remove("like-red");
          }
        }
      } else {
        throw new Error("not ok!");
      }
    } catch (error) {
      console.log(error);
    }
  });
});
