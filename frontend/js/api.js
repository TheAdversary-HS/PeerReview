const prefix = document.getElementsByTagName("base")[0].href.replace(/(?=.*)\/$/gm, "");
async function unknownResponse(resp) {
    let text = await resp.text();
    try {
        let json = JSON.parse(text);
        return new Error(`${json["message"]} (${resp.status})`);
    }
    catch (error) {
        return new Error(`Server sent unknown error: ${text} (${resp.status})`);
    }
}
function buildQuery(options) {
    let query = [];
    options.forEach(element => {
        if (element[1] !== undefined) {
            if (element[1] instanceof Date) {
                element[1] = element[1].getSeconds();
            }
            else if (element[1] instanceof Array) {
                element[1] = JSON.stringify(element[1]);
            }
            query.push(`${element[0]}=${element[1]}`);
        }
    });
    return query.join("&");
}
async function login(username, password) {
    let result = await fetch(`${prefix}/api/login`, {
        method: "POST",
        body: JSON.stringify({ "username": username, "password": password }),
        credentials: "same-origin"
    });
    switch (result.status) {
        case 200:
            return;
        case 401:
            throw new Error("Wrong username and/or password");
        default:
            throw await unknownResponse(result);
    }
}
async function authors(name, limit) {
    let query = [["name", name], ["limit", limit]];
    let result = await fetch(`${prefix}/api/authors?${buildQuery(query)}`);
    if (result.status == 200) {
        return await result.json();
    }
    else {
        throw await unknownResponse(result);
    }
}
async function tags(name, limit) {
    let query = [["name", name], ["limit", limit]];
    let result = await fetch(`${prefix}/api/tags?${buildQuery(query)}`);
    if (result.status == 200) {
        return await result.json();
    }
    else {
        throw await unknownResponse(result);
    }
}
async function recent(limit = 20) {
    let query = [["limit", limit]];
    let result = await fetch(`${prefix}/api/recent?${buildQuery(query)}`);
    if (result.status == 200) {
        return await result.json();
    }
    else {
        throw await unknownResponse(result);
    }
}
async function search(q) {
    let query = [["q", q.query], ["from", q.from], ["to", q.to], ["authors", q.authors], ["tags", q.tags], ["limit", q.limit]];
    let result = await fetch(`${prefix}/api/search?${buildQuery(query)}`);
    if (result.status == 200) {
        return await result.json();
    }
    else {
        throw unknownResponse(result);
    }
}
async function getArticle(id) {
    let query = [["id", id]];
    let result = await fetch(`${prefix}/api/article?${buildQuery(query)}`, {
        method: "GET"
    });
    switch (result.status) {
        case 200:
            return await result.json();
        case 401:
            throw new Error("Not authorized");
        case 404:
            throw new Error("Article not found");
        default:
            throw await unknownResponse(result);
    }
}
async function uploadArticle(payload) {
    let result = await fetch(`${prefix}/api/article`, {
        method: "POST",
        body: JSON.stringify(payload)
    });
    switch (result.status) {
        case 201:
            return await result.json();
        case 401:
            throw new Error("Not authorized");
        case 409:
            throw new Error("An article with the same title already exists");
        default:
            throw await unknownResponse(result);
    }
}
async function editArticle(payload) {
    let result = await fetch(`${prefix}/api/article`, {
        method: "PATCH",
        body: JSON.stringify(payload)
    });
    let json = await result.json();
    switch (result.status) {
        case 201:
            return json;
        case 401:
            throw new Error("Not authorized");
        case 404:
            throw new Error("Could not find article");
        case 409:
            throw new Error("An article with the same title already exists");
        default:
            throw await unknownResponse(result);
    }
}
async function deleteArticle(id) {
    let result = await fetch(`${prefix}/api/article`, {
        method: "DELETE",
        body: JSON.stringify({ "id": id })
    });
    switch (result.status) {
        case 200:
            return;
        case 401:
            throw new Error("Not authorized");
        case 404:
            throw new Error("Could not find article");
        default:
            throw await unknownResponse(result);
    }
}
async function getAssets(name, limit = 20) {
    let query = [["q", name], ["limit", limit]];
    let result = await fetch(`${prefix}/api/assets?${buildQuery(query)}`, {
        method: "GET",
    });
    switch (result.status) {
        case 200:
            return await result.json();
        case 401:
            throw new Error("Not authorized");
        default:
            throw await unknownResponse(result);
    }
}
async function addAsset(name, content) {
    let result = await fetch(`${prefix}/api/assets`, {
        method: "POST",
        body: JSON.stringify({
            "name": name,
            "content": content,
        })
    });
    switch (result.status) {
        case 201:
            return await result.json();
        case 401:
            throw new Error("Not authorized");
        case 409:
            throw new Error("An asset with the same name already exists");
        default:
            throw await unknownResponse(result);
    }
}
async function deleteAsset(id) {
    let result = await fetch(`${prefix}/api/assets`, {
        method: "DELETE",
        body: JSON.stringify({ "id": id })
    });
    switch (result.status) {
        case 200:
            return;
        case 401:
            throw new Error("Not authorized");
        case 404:
            throw new Error("An asset with this id does not exist");
        default:
            throw await unknownResponse(result);
    }
}
//# sourceMappingURL=api.js.map