"use client";
import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { analysis, Result } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";

function ScoreRing({ score }: { score: number }) {
  const color =
    score >= 80 ? "text-green-400" : score >= 60 ? "text-yellow-400" : score >= 40 ? "text-orange-400" : "text-red-400";
  const bg =
    score >= 80 ? "border-green-500" : score >= 60 ? "border-yellow-500" : score >= 40 ? "border-orange-500" : "border-red-500";
  const label =
    score >= 80 ? "Strong Match" : score >= 60 ? "Good Match" : score >= 40 ? "Partial Match" : "Weak Match";

  return (
    <div className="flex flex-col items-center gap-3">
      <div className={`w-36 h-36 rounded-full border-8 ${bg} flex flex-col items-center justify-center`}>
        <span className={`text-4xl font-bold ${color}`}>{score}</span>
        <span className="text-gray-400 text-xs">/ 100</span>
      </div>
      <span className={`text-sm font-semibold ${color}`}>{label}</span>
    </div>
  );
}

export default function ResultPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;
  const [result, setResult] = useState<Result | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!isLoggedIn()) { router.replace("/login"); return; }
    fetchResult();
  }, [id, router]);

  async function fetchResult() {
    try {
      const { data } = await analysis.getResult(id);
      setResult(data);
    } catch {
      setError("Result not found.");
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-4 animate-spin">⚙️</div>
          <p className="text-gray-400">Loading results...</p>
        </div>
      </div>
    );
  }

  if (error || !result) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-4">❌</div>
          <p className="text-red-400">{error}</p>
          <button onClick={() => router.push("/dashboard")} className="mt-4 text-violet-400 hover:underline">
            Back to dashboard
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950">
      {/* Nav */}
      <nav className="border-b border-gray-800 px-6 py-4 flex items-center justify-between">
        <button onClick={() => router.push("/dashboard")} className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors">
          ← Back to Dashboard
        </button>
        <span className="text-gray-500 text-sm">Analysis Result</span>
      </nav>

      <div className="max-w-4xl mx-auto px-6 py-8 space-y-6">
        {/* Header */}
        <div className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
          <div className="flex flex-col md:flex-row items-center gap-8">
            <ScoreRing score={result.score} />
            <div className="flex-1">
              <h1 className="text-2xl font-bold text-white mb-1">Match Analysis</h1>
              <p className="text-gray-400 mb-3">
                <span className="text-violet-400">{result.resume_file_name}</span>
                {" "}&rarr;{" "}
                <span className="text-violet-400">{result.job_title}</span>
              </p>
              {/* Score bar */}
              <div className="w-full bg-gray-800 rounded-full h-3 mb-2">
                <div
                  className={`h-3 rounded-full transition-all duration-1000 ${
                    result.score >= 80 ? "bg-green-500" : result.score >= 60 ? "bg-yellow-500" : result.score >= 40 ? "bg-orange-500" : "bg-red-500"
                  }`}
                  style={{ width: `${result.score}%` }}
                />
              </div>
              <p className="text-gray-500 text-sm">{result.score}% match with job requirements</p>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Strengths */}
          <div className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
            <h2 className="text-lg font-semibold text-green-400 mb-4 flex items-center gap-2">
              ✅ Strengths
            </h2>
            {result.strengths && result.strengths.length > 0 ? (
              <ul className="space-y-2">
                {result.strengths.map((s, i) => (
                  <li key={i} className="flex items-start gap-2 text-gray-300 text-sm">
                    <span className="text-green-500 mt-0.5">•</span>
                    {s}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-gray-500 text-sm">No specific strengths identified.</p>
            )}
          </div>

          {/* Missing Skills */}
          <div className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
            <h2 className="text-lg font-semibold text-red-400 mb-4 flex items-center gap-2">
              ⚠️ Missing Skills
            </h2>
            {result.missing_skills && result.missing_skills.length > 0 ? (
              <ul className="space-y-2">
                {result.missing_skills.map((s, i) => (
                  <li key={i} className="flex items-start gap-2 text-gray-300 text-sm">
                    <span className="text-red-500 mt-0.5">•</span>
                    {s}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-gray-500 text-sm">No significant skill gaps found.</p>
            )}
          </div>
        </div>

        {/* Recommendation */}
        <div className="bg-gray-900 rounded-2xl p-6 border border-gray-800">
          <h2 className="text-lg font-semibold text-violet-400 mb-4">💡 AI Recommendation</h2>
          <p className="text-gray-300 leading-relaxed">{result.recommendation}</p>
        </div>

        {/* Actions */}
        <div className="flex gap-4">
          <button
            onClick={() => router.push("/dashboard")}
            className="flex-1 bg-gray-800 hover:bg-gray-700 text-white font-medium py-3 rounded-xl transition-colors"
          >
            Analyze Another
          </button>
          <button
            onClick={() => router.push(`/ranking/${result.job_id}`)}
            className="flex-1 bg-violet-600 hover:bg-violet-700 text-white font-medium py-3 rounded-xl transition-colors"
          >
            View All Candidates
          </button>
        </div>
      </div>
    </div>
  );
}
