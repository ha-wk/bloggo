package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ha-wk/bloggo/apis/auth"
	"github.com/ha-wk/bloggo/apis/models"
	"github.com/ha-wk/bloggo/apis/responses"
	"github.com/ha-wk/bloggo/apis/utils/formaterror"
)

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := post.SavePost(server.DB) //pata kro kya error he
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)

}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}

	postReceived, err := post.FindPostById(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) //var to store id

	pid, err := strconv.ParseInt(vars["id"], 10, 64) //parsing the id from request body
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r) //extracting user id from body for authentication
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}

	post := models.Post{} //check if the post that we want to update exits or not
	err = server.DB.Debug().Model(post).Where("id=?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	if uid != post.AuthorID { //check if the user attempting to update a post belonging to him or not
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}

	body, err := ioutil.ReadAll(r.Body) //reading the data posted
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate := models.Post{}             //main processing of data
	err = json.Unmarshal(body, &postUpdate) //unmarshalling the data of body
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != postUpdate.AuthorID { //proceed only when requested user id equals the one gotten from token
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}

	postUpdate.Prepare()        //prepare() from models/post.go
	err = postUpdate.Validate() //validate() from models/post.go
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID //it is important to tell model to update post id,other update field are set above

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, postUpdated)

}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//is a valid post?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//is user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}

	post := models.Post{} //check if the post that we want to update exits or not
	err = server.DB.Debug().Model(post).Where("id=?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}

	if uid != post.AuthorID { //check if the user attempting to update a post belonging to him or not
		responses.ERROR(w, http.StatusUnauthorized, errors.New("UNAUTHORIZED"))
		return
	}

	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")

}
