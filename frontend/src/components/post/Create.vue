<template>
   <div class="container">
        <div class="col-md-8 form-wrapper">
          <h2> Добавить запись </h2>
          <form id="create-post-form" @submit.prevent="createPost">
               <div class="form-group col-md-12">
                <label for="title"> Title </label>
                <input type="text" id="title" v-model="title" name="title" class="form-control" placeholder="Enter title">
               </div>
              <div class="form-group col-md-12">
                  <label for="description"> Description </label>
                  <input type="text" id="description" v-model="description" name="description" class="form-control" placeholder="Enter Description">
              </div>
              <div class="form-group col-md-12">
                  <label for="body"> Текст </label>
                  <textarea id="body" cols="30" rows="5" v-model="body" class="form-control"></textarea>
              </div>
              <div class="form-group col-md-12">
                  <label for="author"> Автор </label>
                  <input type="text" id="author" v-model="author" name="author" class="form-control">
              </div>
              <div class="form-group col-md-12">
                  <label for="author"> Категория </label>
                  <select v-model="categoryId" class="form-control">
                    <option v-for="category in categories" v-bind:key="category.id" v-bind:value="category.id">
                      {{ category.title }}
                    </option>
                  </select>
              </div>
              

              <div class="form-group col-md-4 pull-right">
                  <button class="btn btn-success" type="submit"> Добавить </button>
              </div>          
          </form>
        </div>
    </div>
</template>

<script>
import axios from "axios";
import { server } from "../../utils/helper";
import router from "../../router";
export default {
  data() {
    return {
      title: "",
      description: "",
      body: "",
      author: "Example Examplovch",
      date_posted: "",
      categoryId: "",
      categories: this.fetchCategories(),
    };
  },
  created() {
    this.date_posted = new Date().toLocaleDateString();
  },
  methods: {
    fetchCategories() {
      axios
        .get(`${server.baseURL}/blog/categories`)
        .then(data => (this.categories = data.data));
    },
    createPost() {
      let postData = {
        title: this.title,
        description: this.description,
        body: this.body,
        author: this.author,
        categoryId: this.categoryId,
        date_posted: this.date_posted
      };
      this.__submitToServer(postData);
    },
    __submitToServer(data) {
      axios.post(`${server.baseURL}/blog/post`, data).then(data => {
        console.log(data);
        router.push({ name: "home" });
      });
    }
  }
};
</script>