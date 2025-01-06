<script lang="ts">
  import { onMount } from 'svelte';
  import { writable, get } from 'svelte/store';

  let videoElement;
  let videoSelect;

  let deviceInfos = writable<MediaDeviceInfo[]>([]);
  let stream = writable<MediaStream>();

  onMount(() => {
    getStream().then(getDevices).then(gotDevices);
  });

  function getDevices() {
    return navigator.mediaDevices.enumerateDevices();
  }

  function gotDevices(p_deviceInfos) {
    deviceInfos.set(p_deviceInfos);
  }

  function getStream() {
    if (get(stream)) {
      get(stream).getTracks().forEach(track => {
        track.stop();
      });
    }

    const videoSource = videoSelect.value;
    const constraints = {
      video: {deviceId: videoSource ? {exact: videoSource} : undefined}
    };

    return navigator.mediaDevices.getUserMedia(constraints).
      then(gotStream).catch(handleError);
  }

  function gotStream(p_stream) {
    stream.set(p_stream);
    videoSelect.selectedIndex = [...videoSelect.options].
      findIndex(option => option.text === p_stream.getVideoTracks()[0].label);
    videoElement.srcObject = p_stream;
  }

  function handleError(error) {
    console.error('Error: ', error);
  }
</script>

<style>
  main {
    height: 100vh;
  }

  video {
    width: 100%;
    height: 100%;
    background-color: rgb(2, 12, 24);
    object-fit: contain;
  }

  .content {
    height: calc(100% - 40px);
  }

  .footer {
    height: 40px;
  }
</style>

<main>
  <div class="content">
    <video autoplay bind:this={videoElement}></video>
  </div>
  <div class="footer">
    <select bind:this={videoSelect} on:change={getStream}>
      {#each $deviceInfos as deviceInfo, index}
        {#if deviceInfo.kind === 'videoinput'}
          <option value={deviceInfo.deviceId}>{deviceInfo.label || `Camera ${index + 1}`}</option>
        {/if}
      {/each}
    </select>
  </div>
</main>