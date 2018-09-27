<template>
  <div class="hello">


    <audio controls id="speaker" hidden>
      <source src="text_1.mp3" type="audio/mpeg">
    </audio>

    <h1>{{ story.name }}</h1>
    <h3>{{ story.description }}</h3>
    <h3>Let's Start?</h3>

    <!-- <div class="onoffswitch">
        <input type="checkbox" name="onoffswitch" class="onoffswitch-checkbox" id="myonoffswitch" checked>
        <label class="onoffswitch-label" for="myonoffswitch">
            <span class="onoffswitch-inner"></span>
            <span class="onoffswitch-switch"></span>
        </label>
    </div> -->

    <p class="speaker-laura finish">For guide and recipes on how to configure customize this project, check out the other projects story</p>
    <p class="speaker-you current">Hey, my name is Daniel... Im from Brazil!  And you, where are you from?</p>
    <p class="speaker-laura">Im from New Zeland. Its good to talk with you</p>

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
      story: {
        name: "Apresentation Intro",
        description: "You are talking with Laura, She's from New Zeland and want to met new people.",
        script:[
          {
            "text":"Its a night here, and whats your name?",
            "src":"text_laura_1.mp3",
            "speaker":"bot_woman"
          },
          {
            "text":"Hi, my name is Daniel, and you?",
            "result":"Hi, my name is Danielle, and you?",
            "src":"text_you_1.mp3",
            "speaker":"me",
            "duration": 4
          }
        ]
      }
    };
  },
  mounted() {
    var speaker = document.getElementById("speaker");
    //speaker.play();

    speaker.onended = function() {
      alert("The audio has ended");
    };
  },
  created: function() {
    //var vid = document.getElementById("speaker");

    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      console.log("getUserMedia supported.");
      navigator.mediaDevices
        .getUserMedia(
          // constraints - only audio needed for this app
          {
            audio: true
          }
        )

        // Success callback
        .then(stream => {
          var mediaRecorder = new MediaRecorder(stream);
          mediaRecorder.start();
          console.log(mediaRecorder.state);

          setTimeout(() => {
            mediaRecorder.stop();
            console.log(mediaRecorder.state);
          }, 3000);

          var chunks = [];
          mediaRecorder.ondataavailable = function(e) {
            console.log(e);
            chunks.push(e.data);
          };

          mediaRecorder.onstop = function(e) {
            console.log("recorder stopped");
            var blob = new Blob(chunks, { type: "audio/ogg; codecs=opus" });
            console.log(blob);
          };
        })

        // Error callback
        .catch(function(err) {
          console.log("The following getUserMedia error occured: " + err);
        });
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
