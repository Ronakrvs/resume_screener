"use client";
import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useDropzone } from "react-dropzone";
import { resumes, jobs, analysis, Resume, Job } from "@/lib/api";
import { removeToken, isLoggedIn } from "@/lib/auth";

export default function DashboardPage() {
  const router = useRouter();
  const [resumeList, setResumeList] = useState<Resume[]>([]);
  const [jobList, setJobList] = useState<Job[]>([]);
  const [selectedResume, setSelectedResume] = useState("");
  const [selectedJob, setSelectedJob] = useState("");
  const [jobTitle, setJobTitle] = useState("");
  const [jobDescription, setJobDescription] = useState("");
  const [uploading, setUploading] = useState(false);
  const [analyzing, setAnalyzing] = useState(false);
  const [creatingJob, setCreatingJob] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    if (!isLoggedIn()) { router.replace("/login"); return; }
    fetchData();
  }, [router]);

  async function fetchData() {
    try {
      const [rRes, jRes] = await Promise.all([resumes.list(), jobs.list()]);
      setResumeList(rRes.data.resumes || []);
      setJobList(jRes.data.jobs || []);
    } catch {
      setError("Failed to load data. Please try again.");
    }
  }

  const onDrop = useCallback(async (acceptedFiles: File[]) => {
    if (!acceptedFiles.length) return;
    const file = acceptedFiles[0];
    setUploading(true);
    setError("");
    try {
      const { data } = await resumes.upload(file);
      setResumeList((prev) => [data, ...prev]);
      setSelectedResume(data.id);
      setSuccess(`"${file.name}" uploaded successfully`);
      setTimeout(() => setSuccess(""), 3000);
    } catch {
      setError("Failed to upload resume. Check file type (PDF/DOCX/TXT) and size (max 10MB).");
    } finally {
      setUploading(false);
    }
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { "application/pdf": [".pdf"], "application/vnd.openxmlformats-officedocument.wordprocessingml.document": [".docx"], "text/plain": [".txt"] },
    multiple: false,
    disabled: uploading,
  });

  async function handleCreateJob(e: React.FormEvent) {
    e.preventDefault();
    if (!jobTitle.trim() || !jobDescription.trim()) return;
    setCreatingJob(true);
    setError("");
    try {
      const { data } = await jobs.create(jobTitle, jobDescription);
      setJobList((prev) => [data, ...prev]);
      setSelectedJob(data.id);
      setJobTitle("");
      setJobDescription("");
      setSuccess("Job created successfully");
      setTimeout(() => setSuccess(""), 3000);
    } catch {
      setError("Failed to create job.");
    } finally {
      setCreatingJob(false);
    }
  }

  async function handleAnalyze() {
    if (!selectedResume || !selectedJob) {
      setError("Please select both a resume and a job to analyze.");
      return;
    }
    setAnalyzing(true);
    setError("");
    try {
      const { data } = await analysis.analyze(selectedResume, selectedJob);
      router.push(`/results/${data.id}`);
    } catch {
      setError("Analysis failed. Please try again.");
    } finally {
      setAnalyzing(false);
    }
  }

  function handleLogout() {
    removeToken();
    router.push("/login");
  }

  return (
    <div className="min-h-screen bg-gray-950">
      {/* Nav */}
      <nav className="border-b border-gray-800 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-violet-600 rounded-lg flex items-center justify-center text-white font-bold text-sm">AI</div>
          <span className="text-white font-semibold">Resume Screener</span>
        </div>
        <button onClick={handleLogout} className="text-gray-400 hover:text-white text-sm transition-colors">
          Sign out
        </button>
      </nav>

      <div className="max-w-6xl mx-auto px-6 py-8 grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Left column */}
        <div className="space-y-6">
          {/* Upload */}
          <section className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
            <h2 className="text-lg font-semibold text-white mb-4">Upload Resume</h2>
            <div
              {...getRootProps()}
              className={`border-2 border-dashed rounded-xl p-8 text-center cursor-pointer transition-colors ${
                isDragActive ? "border-violet-500 bg-violet-500/10" : "border-gray-700 hover:border-gray-600"
              } ${uploading ? "opacity-50 cursor-not-allowed" : ""}`}
            >
              <input {...getInputProps()} />
              <div className="text-4xl mb-3">{uploading ? "⏳" : "📄"}</div>
              {uploading ? (
                <p className="text-gray-400">Uploading and extracting text...</p>
              ) : isDragActive ? (
                <p className="text-violet-400">Drop the file here</p>
              ) : (
                <>
                  <p className="text-gray-300 font-medium">Drag & drop or click to upload</p>
                  <p className="text-gray-500 text-sm mt-1">PDF, DOCX, TXT · Max 10MB</p>
                </>
              )}
            </div>

            {resumeList.length > 0 && (
              <div className="mt-4">
                <label className="block text-sm font-medium text-gray-300 mb-2">Select resume for analysis</label>
                <select
                  value={selectedResume}
                  onChange={(e) => setSelectedResume(e.target.value)}
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
                >
                  <option value="">-- Choose a resume --</option>
                  {resumeList.map((r) => (
                    <option key={r.id} value={r.id}>{r.file_name}</option>
                  ))}
                </select>
              </div>
            )}
          </section>

          {/* Create Job */}
          <section className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
            <h2 className="text-lg font-semibold text-white mb-4">Add Job Description</h2>
            <form onSubmit={handleCreateJob} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Job Title</label>
                <input
                  value={jobTitle}
                  onChange={(e) => setJobTitle(e.target.value)}
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-violet-500"
                  placeholder="e.g. Senior Backend Engineer"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1">Description</label>
                <textarea
                  value={jobDescription}
                  onChange={(e) => setJobDescription(e.target.value)}
                  rows={6}
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-violet-500 resize-none"
                  placeholder="Paste the full job description here..."
                />
              </div>
              <button
                type="submit"
                disabled={creatingJob || !jobTitle.trim() || !jobDescription.trim()}
                className="w-full bg-gray-700 hover:bg-gray-600 disabled:opacity-40 text-white font-medium py-2.5 rounded-lg transition-colors"
              >
                {creatingJob ? "Saving..." : "Save Job"}
              </button>
            </form>
          </section>
        </div>

        {/* Right column */}
        <div className="space-y-6">
          {/* Job selector */}
          {jobList.length > 0 && (
            <section className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
              <h2 className="text-lg font-semibold text-white mb-4">Select Job</h2>
              <select
                value={selectedJob}
                onChange={(e) => setSelectedJob(e.target.value)}
                className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2.5 text-white focus:outline-none focus:ring-2 focus:ring-violet-500"
              >
                <option value="">-- Choose a job --</option>
                {jobList.map((j) => (
                  <option key={j.id} value={j.id}>{j.title}</option>
                ))}
              </select>
              {selectedJob && (
                <button
                  onClick={() => router.push(`/ranking/${selectedJob}`)}
                  className="mt-3 text-sm text-violet-400 hover:text-violet-300 transition-colors"
                >
                  View candidate ranking for this job →
                </button>
              )}
            </section>
          )}

          {/* Analyze CTA */}
          <section className="bg-gradient-to-br from-violet-900/40 to-purple-900/20 rounded-2xl p-6 border border-violet-800/50">
            <h2 className="text-lg font-semibold text-white mb-2">Analyze Match</h2>
            <p className="text-gray-400 text-sm mb-6">
              Select a resume and job description above, then run AI analysis to get a match score, strengths, and skill gaps.
            </p>

            {error && (
              <div className="bg-red-900/40 border border-red-700 text-red-300 px-4 py-3 rounded-lg text-sm mb-4">
                {error}
              </div>
            )}
            {success && (
              <div className="bg-green-900/40 border border-green-700 text-green-300 px-4 py-3 rounded-lg text-sm mb-4">
                {success}
              </div>
            )}

            <div className="flex gap-3 mb-4">
              <div className={`flex-1 rounded-lg p-3 text-sm ${selectedResume ? "bg-violet-800/40 text-violet-300" : "bg-gray-800 text-gray-500"}`}>
                {selectedResume ? `✓ Resume selected` : "No resume selected"}
              </div>
              <div className={`flex-1 rounded-lg p-3 text-sm ${selectedJob ? "bg-violet-800/40 text-violet-300" : "bg-gray-800 text-gray-500"}`}>
                {selectedJob ? `✓ Job selected` : "No job selected"}
              </div>
            </div>

            <button
              onClick={handleAnalyze}
              disabled={analyzing || !selectedResume || !selectedJob}
              className="w-full bg-violet-600 hover:bg-violet-700 disabled:opacity-40 disabled:cursor-not-allowed text-white font-semibold py-3 rounded-xl transition-colors text-lg"
            >
              {analyzing ? "Analyzing with AI..." : "Run AI Analysis"}
            </button>
          </section>

          {/* Stats */}
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-gray-900 rounded-xl p-4 border border-gray-800 text-center">
              <div className="text-3xl font-bold text-violet-400">{resumeList.length}</div>
              <div className="text-gray-400 text-sm mt-1">Resumes</div>
            </div>
            <div className="bg-gray-900 rounded-xl p-4 border border-gray-800 text-center">
              <div className="text-3xl font-bold text-violet-400">{jobList.length}</div>
              <div className="text-gray-400 text-sm mt-1">Jobs</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
