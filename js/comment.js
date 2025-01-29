const AllComments = async (add) => {
  const url = location.pathname;
  try {
    let path = `${url}/comments${location.search}`;
    if (add) {
      path = `${url}/comments`;
    }
    const response = await fetch(path);
    if (!response.ok) {
      throw new Error("err");
    }
    // console.log(response)
    if (response.redirected) {
      location.href = response.url;
    }

    const data = await response.json();
    const comments = data.Comments;
    const DivComments = document.createElement("div");
    const h2 = (document.createElement("h2").textContent = "Comments");
    const h4 = document.createElement("h4");
    h4.classList.add("email-error", "post-error");
    DivComments.append(h2, h4);

    comments.forEach((comment) => {
      const com = document.createElement("div");
      com.className = "comment";

      const userPost = document.createElement("div");
      userPost.className = "user-post";

      const user = document.createElement("div");
      user.className = "a-user-post";

      const img = document.createElement("img");
      img.className = "img-user-post";
      img.src = `/images/pics/${comment.Image}`;

      const h3 = document.createElement("h3");
      h3.textContent = comment.UserName;

      const h6 = document.createElement("h6");
      h6.className = "duration-post";
      h6.textContent = `• ${comment.Duration} ago`;

      user.append(img, h3, h6);
      userPost.append(user);

      const p = document.createElement("p");
      p.id = `p${comment.Id}`;
      p.className = "post-body";
      p.textContent = comment.Body;
      // com.appendChild(p)

      const form = document.createElement("form");
      form.action = `/comments/likes`;
      form.method = `POST`;
      form.addEventListener("submit", (e) => e.preventDefault());

      const div = document.createElement("div");
      div.classList.add("likes-comment");

      const btnLike = document.createElement("button");
      btnLike.name = `like`;
      btnLike.value = `1`;
      btnLike.classList.add("btn-like");

      const spanLike = document.createElement("span");
      spanLike.className = "like-number";
      spanLike.textContent = comment.Like;

      // const svg = document.createElement('svg');

      btnLike.innerHTML = `<svg class="${
        comment.Did ? (comment.Liked ? "like-blue" : "") : ""
      }"
      xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
      fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
      stroke-linejoin="round" class="feather feather-thumbs-up">
      <path
      d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3">
      </path>
      </svg>`;
      btnLike.append(spanLike);

      const btnDisLike = document.createElement("button");
      btnDisLike.name = `like`;
      btnDisLike.value = `0`;
      btnDisLike.classList.add("btn-like");

      const spanDisLike = document.createElement("span");
      spanDisLike.className = "like-number";
      spanDisLike.textContent = comment.Dislike;

      btnDisLike.innerHTML = `
      <svg class="${comment.Did ? (comment.Liked ? "" : "like-red") : ""}"
      xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"
      fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
      stroke-linejoin="round" class="feather feather-thumbs-down">
      <path
      d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zm7-13h2.67A2.31 2.31 0 0 1 22 4v7a2.31 2.31 0 0 1-2.33 2H17">
      </path>
      </svg>
      `;
      btnDisLike.append(spanDisLike);

      div.append(btnLike, btnDisLike);
      // events
      btnLike.addEventListener("click", async (e) => {
        let data = new FormData();
        data.append("like", "1");
        data.append("comment_id", comment.Id);
        let res = await fetch("/comments/likes", {
          method: "POST",
          body: data,
        });

        if (res.redirected) {
          location.href = res.url;
          return;
        }
        const svgLike = btnLike.querySelector("svg");
        const svgDisLike = btnDisLike.querySelector("svg");
        const result = await res.json();
        // e.target.querySelector("span").textContent = result.Likes
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
      });

      btnDisLike.addEventListener("click", async (e) => {
        let data = new FormData();
        data.append("like", "0");
        data.append("comment_id", comment.Id);
        let res = await fetch("/comments/likes", {
          method: "POST",
          body: data,
        });
        if (res.redirected) {
          location.href = res.url;
          return;
        }
        const svgLike = btnLike.querySelector("svg");
        const svgDisLike = btnDisLike.querySelector("svg");
        const result = await res.json();
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
      });

      form.appendChild(div);

      com.append(userPost, p, form);
      DivComments.append(com);
    });

    const chev = document.createElement("div");
    chev.className = "chev";

    //
    const Pahref = document.createElement("a");
    Pahref.href = data.Previous;
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.setAttribute("xmlns", "http://www.w3.org/2000/svg");
    svg.setAttribute("width", "24");
    svg.setAttribute("height", "24");
    svg.setAttribute("viewBox", "0 0 24 24");
    svg.setAttribute("fill", "none");
    svg.setAttribute("stroke", "currentColor");
    svg.setAttribute("stroke-width", "2");
    svg.setAttribute("stroke-linecap", "round");
    svg.setAttribute("stroke-linejoin", "round");
    svg.setAttribute("class", "feather feather-chevron-left");
    const polyline = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "polyline"
    );
    polyline.setAttribute("points", "15 18 9 12 15 6");
    svg.appendChild(polyline);

    //
    const Nahref = document.createElement("a");
    Nahref.href = data.Next;
    const svg2 = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg2.setAttribute("xmlns", "http://www.w3.org/2000/svg2");
    svg2.setAttribute("width", "24");
    svg2.setAttribute("height", "24");
    svg2.setAttribute("viewBox", "0 0 24 24");
    svg2.setAttribute("fill", "none");
    svg2.setAttribute("stroke", "currentColor");
    svg2.setAttribute("stroke-width", "2");
    svg2.setAttribute("stroke-linecap", "round");
    svg2.setAttribute("stroke-linejoin", "round");
    svg2.setAttribute("class", "feather feather-chevron-right");

    const polyline2 = document.createElementNS(
      "http://www.w3.org/2000/svg",
      "polyline"
    );
    polyline2.setAttribute("points", "9 18 15 12 9 6");

    svg2.appendChild(polyline2);
    Pahref.append(svg);
    Nahref.append(svg2);
    if (data.Previous != "0") {
      chev.append(Pahref);
    }
    chev.append(Nahref);
    DivComments.append(chev); // Append the SVG to the body (or any other container)

    svg.addEventListener("click", () => {
      AllComments();
    });

    document.getElementsByClassName("comments")[0].innerHTML = "";
    document.getElementsByClassName("comments")[0].append(DivComments);
    document.querySelector(".post-error").textContent = getCookie("errors");
  } catch (error) {
    console.log(error);
  }
};

// document.querySelector(".comments").addEventListener("load", (e) => {
//   console.log("kk")
// });

AllComments();
function getCookie(name) {
  // Split all cookies into an array
  const cookies = document.cookie.split(";");

  // Loop through the cookies to find the one with the matching name
  for (let cookie of cookies) {
    // Trim any leading or trailing whitespace
    cookie = cookie.trim();

    // Check if the cookie starts with the name we're looking for
    if (cookie.startsWith(name + "=")) {
      // Return the cookie value (everything after the '=')
      return cookie.substring(name.length + 1);
    }
  }

  // Return null if the cookie is not found
  return null;
}

const addCommentForm = document.querySelectorAll(".add-comment-form");
addCommentForm.forEach((comment) => {
  comment.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = new FormData(comment);
    data.append("js", true);
    try {
      const response1 = await fetch("http://localhost:8080/comments/store", {
        method: "POST",
        body: data,
      });
      if (response1.ok) {
        if (response1.redirected) {
          location.href = response1.url;
        } else {
          await AllComments(true);
        }
      } else {
        const resParse = await response1.json();
        if (resParse.Message.startsWith("You are on cooldown")) {
          document.querySelector(".post-error").textContent = resParse.Message;
        } else {
          const divError = document.createElement("div");
          divError.className = "error";
          divError.textContent = `${resParse.Code} | ${resParse.Message}`;
          const main = document.querySelector("main");
          main.insertAdjacentElement("afterend", divError);
          main.insertAdjacentHTML("afterend", "<div></div>");
          main.remove();
        }
      }
    } catch (error) {
      console.log(error);
    }
  });
});
