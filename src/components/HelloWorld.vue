<template>
  <div class="hello">


    <pre class="code">{{story | json}}</pre>

    <h1>{{ story.name }}</h1>
    <h3>{{ story.description }}</h3>
    <h3 v-on:click="start">Let's Start?</h3>

    <!-- <div class="onoffswitch">
        <input type="checkbox" name="onoffswitch" class="onoffswitch-checkbox" id="myonoffswitch" checked>
        <label class="onoffswitch-label" for="myonoffswitch">
            <span class="onoffswitch-inner"></span>
            <span class="onoffswitch-switch"></span>
        </label>
    </div> -->

    <div v-for="script in story.scripts" :key="script.id">
        <p v-bind:class="[script.speaker, script.status]">{{script.text}} {{script.status}}</p>

        <audio controls :id="'speaker_'+script.id" hidden>
          <source :src="script.src" type="audio/mpeg">
        </audio>
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
      //mediaRecorder: MediaRecorder,
      test: 'hello',
      story: {
        name: "Apresentation Intro",
        description: "You are talking with Laura, She's from New Zeland and want to met new people.",
        scripts:[
          {
            "id": "1",
            "text":"It wasn’t just that I was 41, which, let’s face it, isn’t old. It was that I was 41 and bored. And a little tired. And, at times, cantankerous. Crotchety, you might say.",
            "src":"text_laura_1.mp3",
            "speaker":"speaker-laura"
          },
          {
            "id": "2",
            "text":"Exacerbating this problem was the fact that I had spent the entire span of my thirties at one place — a prestigious men’s magazine. I thought I had stability and security and swagger.",
            "result":"Hi, my name is Danielle, and you?",
            "src":"text_you_1.mp3",
            "speaker":"speaker-you",
            "duration": 4000
          },
          {
            "id": "3",
            "text":"What I didn’t realize is that I had slowly started draining energy from the place where I worked instead of injecting it with my own. I was getting soft. I was getting lazy.",
            "src":"text_laura_2.mp3",
            "speaker":"speaker-laura"
          },
          {
            "id": "4",
            "text":"A couple months into unemployment, I got a job at another prestigious men’s magazine. ",
            "src":"text_you_2.mp3",
            "speaker":"speaker-you",
            "duration": 6000,
            "isLast": true,
          },
        ]
      }
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
        this.test = 'opa';
        this.story = this.mergeDeep(this.story, {scripts: {"aa":"sdda"}});
        console.log(this.story.scripts[0].status)
        //console.log(this.story.scripts[0]);
        this.playAudio(0)
      },
      nextStep: function(idx) {
        this.story.scripts[idx-1].status = "finish";
        const script = this.story.scripts[idx];

        if (script === undefined) {
          console.log('finish')
            return;
        }

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
          const id = this.story.scripts[idx].id
          var speaker = document.getElementById("speaker_" + id);
          speaker.play();

          var self = this;
          speaker.onended = function() {
            self.nextStep(idx + 1)
          };
      },
      recordAudio: function(idx){
        const script = this.story.scripts[idx];

         navigator.mediaDevices.getUserMedia({audio: true}).then(stream => {
           var self = this;

            var mediaRecorder = new MediaRecorder(stream);
            mediaRecorder.start();
            console.log(mediaRecorder.state);

            setTimeout(() => {
              if (script.isLast) {
                mediaRecorder.stop();
              }else{
                mediaRecorder.pause();
              }
            }, script.duration);

            var chunks = [];
            mediaRecorder.ondataavailable = function(e) {
              console.log(e);
              chunks.push(e.data);
            };

            mediaRecorder.onstop = function(e) {
              console.log("recorder stopped");
              var blob = new Blob(chunks, { type: "audio/ogg; codecs=opus" });
            };

            mediaRecorder.onpause = function(e) {
              console.log("recorder paused");
              self.nextStep(idx + 1)
            };
          }).catch(function(err) {
            console.log("The following getUserMedia error occured: " + err);
        });
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
  created: function() {

    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      console.log("getUserMedia supported.");
      // navigator.mediaDevices.getUserMedia({audio: true}).then(stream => {
      //   var mediaRecorder = new MediaRecorder(stream);
      //   mediaRecorder.start();
      //   console.log(mediaRecorder.state);

      //   setTimeout(() => {
      //     mediaRecorder.stop();
      //     console.log(mediaRecorder.state);
      //   }, 3000);

      //   var chunks = [];
      //   mediaRecorder.ondataavailable = function(e) {
      //     console.log(e);
      //     chunks.push(e.data);
      //   };

      //   mediaRecorder.onstop = function(e) {
      //     console.log("recorder stopped");
      //     var blob = new Blob(chunks, { type: "audio/ogg; codecs=opus" });
      //     console.log(blob);
      //   };
      // }).catch(function(err) {
      //   console.log("The following getUserMedia error occured: " + err);
      // });
    } else {
      console.log("getUserMedia not supported on your browser!");
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

.finish {
  text-decoration-line:line-through;
}

.current {
  border-left: 5px solid #c3c3c3;
}

.speaker-laura {
  color: #8f7d90;
  font-style: italic;
  
}
.speaker-you {
  background-color: #efead1;
}
h1 {
  font-weight: 700;
  font-size: 32px;
  letter-spacing: -.01em;
  text-decoration-line: underline;
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
