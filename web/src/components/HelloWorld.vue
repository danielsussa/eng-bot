<template>
  <div class="hello" v-if="!loading">
    
    <div class="presentation-container">
        <h1>{{ story.name }}</h1>
        <h3>{{ story.description }}</h3>
        <h3 v-on:click="start">Let's Start?</h3>
    </div>


    <!-- <div class="onoffswitch">
        <input type="checkbox" name="onoffswitch" class="onoffswitch-checkbox" id="myonoffswitch" checked>
        <label class="onoffswitch-label" for="myonoffswitch">
            <span class="onoffswitch-inner"></span>
            <span class="onoffswitch-switch"></span>
        </label>
    </div> -->

    <div class="script-container">
      <div v-for="script in story.scripts" :key="script.id" class="text-container">

          <span>
            <i v-if="script.speaker !== 'speaker-you'" class="fas fa-volume-up text-icon" v-bind:class="{current: script.isCurrent}"></i>
            <i v-if="script.speaker === 'speaker-you'" class="fas fa-microphone text-icon" v-bind:class="{current: script.isCurrent}"></i>
          </span>

          <p v-bind:class="[script.speaker, script.status]">{{script.text}}</p>

          <audio controls :id="'speaker_'+script.id" hidden>
            <source :src="script.src" type="audio/mpeg">
          </audio>
      </div>
    </div>


  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    msg: String
  },
  data() {
    return {
      mediaRecorder: MediaRecorder,
      test: 'hello',
      loading: true,
      story: {}
    };
  },
  mounted() {
    // var speaker = document.getElementById("speaker");
    // //speaker.play();

    // speaker.onended = function() {
    //   alert("The audio has ended");
    // };
  },
  methods: {
      start: function (event) {
        this.nextStep(0)
      },
      nextStep: function(idx) {
        const script = this.story.scripts[idx];

        if (idx !== 0) {
          this.story = this.mergeDeep(this.story, {scripts: {[idx-1]: {status: 'finish', isCurrent: false}}});
        }

        if (script === undefined) {
            return;
        }

        this.story = this.mergeDeep(this.story, {scripts: {[idx]: {isCurrent: true}}});



        //open Mic
        if (script.speaker === "speaker-you") {
            this.recordAudio(idx);
        }

        //Play Sound
        if (script.speaker !== "speaker-you") {
            this.playAudio(idx);
        }

      },
      playAudio: function(idx){
          const id = this.story.scripts[idx].id;
          var speaker = document.getElementById("speaker_" + id);
          speaker.play();

          var self = this;
          speaker.onended = function() {
            self.nextStep(idx + 1)
          };
      },
      recordAudio: function(idx){
        const script = this.story.scripts[idx];

        if (this.mediaRecorder.state === 'paused') {
          this.mediaRecorder.resume();
        } else {
          this.mediaRecorder.start();
        }

        var self = this;
        setTimeout(() => {
          if (script.isLast) {
            console.log('fim!')
            self.mediaRecorder.requestData();

            self.nextStep(idx + 1);
          }else{
            self.mediaRecorder.pause();
            self.nextStep(idx + 1);
          }
        }, script.duration * 1000);

        var chunks = [];
        this.mediaRecorder.ondataavailable = function(e) {
          chunks.push(e.data);
          const audioBlob = new Blob(chunks, { type: "audio/ogg; codecs=opus" });
          var formData = new FormData();
          formData.append('audio', audioBlob);

          var r = confirm("Send to backend")
          if (r == true) {       
            self.$http.post('http://192.168.0.7:1323/story/audio/send/123',formData).
            then(function(data){
                console.log(data);
            });
          }


        };

        // this.mediaRecorder.onstop = function(e) {
        //   console.log("recorder stopped");
        //   var blob = new Blob(chunks, { type: "audio/ogg; codecs=opus" });
        //   self.nextStep(idx + 1)
        // };

        // this.mediaRecorder.onpause = function(e) {
        //   alert("recorder paused");
        //   self.nextStep(idx + 1)
        // };
      },
      mergeDeep: function(target, source) {
          const output = Object.assign({}, target);
          if (this.isObject(target) && this.isObject(source)) {
              Object.keys(source).forEach(key => {
                  if (this.isObject(source[key])) {
                      if (!(key in target))
                          Object.assign(output, { [key]: source[key] });
                      else
                          output[key] = this.mergeDeep(target[key], source[key]);
                  } else {
                      Object.assign(output, { [key]: source[key] });
                  }
              });
          }
          return output;
      },
      isObject: function(item) {
          return (item && typeof item === 'object' && !Array.isArray(item));
      }

  },
  ready: function () {
    console.log('opa')

  },
  startMediaDevice() {

  },
  created: function() {

    this.$http.get('http://192.168.0.7:1323/story/12').then(res => {
      this.story = res.data
      this.loading = false;
    }).catch(function(err) {
      alert(err.message);
    })

    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
        // Get media

        navigator.mediaDevices.getUserMedia({audio: true, video: false}).then(stream => {
          var options = {
            //mimeType : 'audio/ogg',
            audioBitsPerSecond : 16000,
          }
          this.mediaRecorder = new MediaRecorder(stream, options);

        }).catch(function(err) {
            alert(err)
            console.log("The following getUserMedia error occured: " + err);
        });
    } else {
      alert("getUserMedia not supported on your browser!")
      console.log("getUserMedia not supported on your browser!");
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

.presentation-container {
  //background-color: black;
}

.script-container {
  /* -webkit-mask-image: -webkit-linear-gradient(rgba(0,0,0,0), rgba(0,0,0,1), rgba(0,0,0,1),rgba(0,0,0,1),rgba(0,0,0,1)); */
}

.text-container {
  display: flex;
  align-items: center;
  justify-content: center;

  span {
    .text-icon {
      padding-top: 15px;
      margin-right: 10px;
      width: 20px;
      height: 30px;
      font-size: 13px;
      border-radius: 20px;
      color: #bfa2a2;
      background-color: #efebd1;
      border-left: 3px solid #c3c3c3;

      -webkit-transition: background-color 0.3s ease-out;
      -moz-transition: background-color 0.3s ease-out;
      -o-transition: background-color 0.3s ease-out;
      transition: background-color 0.3s ease-out;

      &.current {
          -webkit-transition: background-color 0.3s ease-out;
          -moz-transition: background-color 0.3s ease-out;
          -o-transition: background-color 0.3s ease-out;
          transition: background-color 0.3s ease-out;
          background-color: #ffc9c9;
      }
    }
  }

}



.current {
  border-left: 3px solid #c3c3c3;
}

.speaker-laura {
  color: #8f7d90;
  font-style: italic;
  
}
.speaker-you {
  background-color: #f1efe7;
  color: #6d6a5c;
}
h1 {
  font-weight: 700;
  font-size: 32px;
  letter-spacing: -.01em;
  text-decoration-line: underline;
  color:#636363;
  margin-bottom: 0px;
}
h3 {
  font-weight: 400;
  font-size: 16px;
  color: #757575;
  max-width: 420px;
  margin: auto;
  letter-spacing: -.01em;
}
p {
  font-weight: 400;
  font-size: 19px;
  letter-spacing: -.01em;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

@media only screen and (max-width: 480px) {
    p {
      font-size: 16px;
    }
}

.onoffswitch {

   margin: auto;
   margin-top: 20px;
    position: relative; width: 200px;
    -webkit-user-select:none; -moz-user-select:none; -ms-user-select: none;
}
.onoffswitch-checkbox {
    display: none;
}
.onoffswitch-label {
    display: block; overflow: hidden; cursor: pointer;
    border: 2px solid #FFFFFF; border-radius: 7px;
}
.onoffswitch-inner {
    display: block; width: 200%; margin-left: -100%;
    transition: margin 0.3s ease-in 0s;
}
.onoffswitch-inner:before, .onoffswitch-inner:after {
    display: block; float: left; width: 50%; height: 25px; padding: 0; line-height: 25px;
    font-size: 14px; color: white; font-family: Trebuchet, Arial, sans-serif; font-weight: bold;
    box-sizing: border-box;
}
.onoffswitch-inner:before {
    content: "Let's go!";
    font-family: 'Noto Serif KR', sans-serif;
    padding-left: 0px;
    background-color: #A0B9BF; color: #FFFFFF;
}
.onoffswitch-inner:after {
    content: "Let's train!";
    font-family: 'Noto Serif KR', sans-serif;
    padding-right: 60px;
    background-color: #CCB7C9; color: #757575;
    text-align: right;
}
.onoffswitch-switch {
    display: block; width: 8px; margin: 8.5px;
    background: #FFFFFF;
    position: absolute; top: 0; bottom: 0;
    right: 171px;
    border: 2px solid #FFFFFF; border-radius: 7px;
    transition: all 0.3s ease-in 0s; 
}
.onoffswitch-checkbox:checked + .onoffswitch-label .onoffswitch-inner {
    margin-left: 0;
}
.onoffswitch-checkbox:checked + .onoffswitch-label .onoffswitch-switch {
    right: 0px; 
}
</style>
