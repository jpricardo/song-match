export function encodeWAV(audioBuffer: AudioBuffer): ArrayBuffer {
	const sampleRate = audioBuffer.sampleRate;
	const channelData = audioBuffer.getChannelData(0); // mono
	const length = channelData.length;

	// Create WAV file structure
	const WAV_HEADER_SIZE = 44;
	const buffer = new ArrayBuffer(WAV_HEADER_SIZE + length * 2);
	const view = new DataView(buffer);

	// WAV header
	const writeString = (offset: number, string: string) => {
		for (let i = 0; i < string.length; i++) {
			view.setUint8(offset + i, string.charCodeAt(i));
		}
	};

	writeString(0, 'RIFF');
	view.setUint32(4, 36 + length * 2, true);
	writeString(8, 'WAVE');
	writeString(12, 'fmt ');
	view.setUint32(16, 16, true); // fmt chunk size
	view.setUint16(20, 1, true); // PCM
	view.setUint16(22, 1, true); // mono
	view.setUint32(24, sampleRate, true);
	view.setUint32(28, sampleRate * 2, true); // byte rate
	view.setUint16(32, 2, true); // block align
	view.setUint16(34, 16, true); // bits per sample
	writeString(36, 'data');
	view.setUint32(40, length * 2, true);

	// PCM data (convert float [-1, 1] to int16)
	let offset = 44;
	for (let i = 0; i < length; i++) {
		const sample = Math.max(-1, Math.min(1, channelData[i]));
		view.setInt16(offset, sample < 0 ? sample * 0x8000 : sample * 0x7fff, true);
		offset += 2;
	}

	return buffer;
}
