import { useEffect, useState } from "react"
import { SendUrl } from "../wailsjs/go/main/App"
import { EventsOn } from "../wailsjs/runtime"

function removeOldDuplicates(dataList: any) {
	const uniqueDataMap = new Map()

	dataList.forEach((data: any) => {
		const existingData = uniqueDataMap.get(data.name)

		if (!existingData || data.date > existingData.date) {
			uniqueDataMap.set(data.name, data)
		}
	})

	const uniqueDataList = Array.from(uniqueDataMap.values())

	return uniqueDataList.sort((a, b) => b.progressInNumber - a.progressInNumber)
}

export default function App() {
	const [url, setUrl] = useState("")
	const [progress, setProgress] = useState(0)
	const [imagesProgress, setImagesProgress] = useState<any[]>([])

	const [logMessages, setLogMessages] = useState<string[]>([])
	const [galleryName, setGalleryName] = useState<string>(" ")
	const [infoTask, setInfoTask] = useState(" ")
	const [status, setStatus] = useState<"started" | "stopped">("stopped")

	useEffect(() => {
		EventsOn("log_event", (message) => {
			setLogMessages((prevMessages) => [message, ...prevMessages])
		})
	}, [])

	useEffect(() => {
		EventsOn("gallery_name", (message) => setGalleryName(message))
	}, [])

	useEffect(() => {
		EventsOn("info_task", (message) => setInfoTask(message))
	}, [])

	useEffect(() => {
		EventsOn("progress", (p) => setProgress(+p))
	}, [])

	useEffect(() => {
		EventsOn("image_progress", (p) => {
			const newProgress = JSON.parse(p || "")
			setImagesProgress((prev) => [
				...prev,
				{ ...newProgress, date: new Date(), progressInNumber: parseFloat(newProgress.progress.replace("%", "")) },
			])
		})
	}, [])

	useEffect(() => {
		EventsOn("status", (s) => setStatus(s))
	}, [])

	const handleSendUrl = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		await SendUrl(url)
	}

	return (
		<>
			<main>
				{status !== "started" && (
					<>
						<form onSubmit={handleSendUrl} className="url-section">
							<fieldset className="cyber-input ac-purple">
								<input
									type="text"
									name="url"
									placeholder="Enter url..."
									value={url}
									onChange={(e) => setUrl(e.target.value)}
									autoComplete="off"
								/>
							</fieldset>

							<button type="submit" className="cyber-button-small bg-purple fg-white">
								Enter
								<span className="glitchtext">Scrap</span>
							</button>
						</form>

						<div className="message-info">{infoTask}</div>
					</>
				)}

				<section className="progress-info">
					<div className="">{galleryName}</div>

					{status === "started" && !!progress && (
						<div className="progress-bar">
							<div style={{ width: `${progress}%` }} />
						</div>
					)}

					<div className="message-info">{infoTask}</div>

					<div className="image-download-progress">
						<ul>
							{removeOldDuplicates(imagesProgress).map((progress) => (
								<li key={progress.name}>
									<div className="image-download-stats">
										<div>{progress.name}</div>
										<div>{progress.progress}</div>
									</div>

									<div className="img-progress">
										<div style={{ width: progress.progress }} className="img-progress-bar" />
									</div>

									<div className="image-propriety">
										{progress.written} / {progress.total}
									</div>
								</li>
							))}
						</ul>
					</div>
				</section>

				<section className="entries-section">
					{logMessages.map((message, index) => (
						<div key={message + index}>{message}</div>
					))}
				</section>
			</main>
		</>
	)
}
