'use client';

export function useMediaStream() {
	return navigator.mediaDevices.getUserMedia({ audio: true });
}
