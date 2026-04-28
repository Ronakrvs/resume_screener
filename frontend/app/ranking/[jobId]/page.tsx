"use client";
import { useEffect, useState } from "react";
import { useRouter, useParams } from "next/navigation";
import { analysis, Result } from "@/lib/api";
import { isLoggedIn } from "@/lib/auth";

function ScoreBadge({ score }: { score: number }) {
  const cls =
    score >= 80
      ? "bg-green-900/40 text-green-400 border-green-700"
      : score >= 60
      ? "bg-yellow-900/40 text-yellow-400 border-yellow-700"
      : score >= 40
      ? "bg-orange-900/40 text-orange-400 border-orange-700"
      : "bg-red-900/40 text-red-400 border-red-700";
  return (
    <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-bold border ${cls}`}>
      {score}
    </span>
  );
}

export default function RankingPage() {
  const router = useRouter();
  const params = useParams();
  const jobId = params.jobId as string;
  const [results, setResults] = useState<Result[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filterMin, setFilterMin] = useState(0);

  useEffect(() => {
    if (!isLoggedIn()) { router.replace("/login"); return; }
    fetchRanking();
  }, [jobId, router]);

  async function fetchRanking() {
    try {
      const { data } = await analysis.getRanking(jobId);
      setResults(data.results || []);
    } catch {
      setError("Failed to load ranking.");
    } finally {
      setLoading(false);
    }
  }

  const filtered = results.filter((r) => r.score >= filterMin);
  const jobTitle = results[0]?.job_title ?? "Job";

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-400">Loading ranking...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-950">
      <nav className="border-b border-gray-800 px-6 py-4 flex items-center justify-between">
        <button onClick={() => router.push("/dashboard")} className="text-gray-400 hover:text-white transition-colors">
          ← Back to Dashboard
        </button>
        <span className="text-gray-500 text-sm">Candidate Ranking</span>
      </nav>

      <div className="max-w-5xl mx-auto px-6 py-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-white">Candidate Ranking</h1>
          <p className="text-gray-400 mt-1">Job: <span className="text-violet-400">{jobTitle}</span></p>
        </div>

        {error && (
          <div className="bg-red-900/40 border border-red-700 text-red-300 px-4 py-3 rounded-lg text-sm mb-6">
            {error}
          </div>
        )}

        {/* Filters */}
        <div className="bg-gray-900 rounded-xl p-4 border border-gray-800 mb-6 flex flex-wrap gap-4 items-center">
          <span className="text-gray-400 text-sm">Min score:</span>
          {[0, 40, 60, 80].map((v) => (
            <button
              key={v}
              onClick={() => setFilterMin(v)}
              className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
                filterMin === v ? "bg-violet-600 text-white" : "bg-gray-800 text-gray-400 hover:bg-gray-700"
              }`}
            >
              {v === 0 ? "All" : `${v}+`}
            </button>
          ))}
          <span className="text-gray-500 text-sm ml-auto">{filtered.length} candidates</span>
        </div>

        {/* Table */}
        {filtered.length === 0 ? (
          <div className="text-center py-16 text-gray-500">
            {results.length === 0
              ? "No candidates analyzed for this job yet."
              : `No candidates with score ≥ ${filterMin}.`}
          </div>
        ) : (
          <div className="bg-gray-900 rounded-2xl border border-gray-800 overflow-hidden">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-800">
                  <th className="text-left px-6 py-4 text-gray-400 text-sm font-medium">Rank</th>
                  <th className="text-left px-6 py-4 text-gray-400 text-sm font-medium">Resume</th>
                  <th className="text-left px-6 py-4 text-gray-400 text-sm font-medium">Score</th>
                  <th className="text-left px-6 py-4 text-gray-400 text-sm font-medium">Top Strengths</th>
                  <th className="text-left px-6 py-4 text-gray-400 text-sm font-medium">Gaps</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((r, idx) => (
                  <tr key={r.id} className="border-b border-gray-800/50 hover:bg-gray-800/30 transition-colors">
                    <td className="px-6 py-4">
                      <span className={`text-lg font-bold ${idx === 0 ? "text-yellow-400" : idx === 1 ? "text-gray-300" : idx === 2 ? "text-amber-600" : "text-gray-500"}`}>
                        #{idx + 1}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-white text-sm font-medium">{r.resume_file_name}</span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <ScoreBadge score={r.score} />
                        <div className="w-24 bg-gray-800 rounded-full h-1.5">
                          <div
                            className={`h-1.5 rounded-full ${r.score >= 80 ? "bg-green-500" : r.score >= 60 ? "bg-yellow-500" : r.score >= 40 ? "bg-orange-500" : "bg-red-500"}`}
                            style={{ width: `${r.score}%` }}
                          />
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-1">
                        {(r.strengths || []).slice(0, 2).map((s, i) => (
                          <span key={i} className="bg-green-900/30 text-green-400 text-xs px-2 py-0.5 rounded-full border border-green-800/50">
                            {s.length > 25 ? s.slice(0, 25) + "…" : s}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="text-red-400 text-xs">
                        {(r.missing_skills || []).length} gap{(r.missing_skills || []).length !== 1 ? "s" : ""}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => router.push(`/results/${r.id}`)}
                        className="text-violet-400 hover:text-violet-300 text-sm transition-colors whitespace-nowrap"
                      >
                        View →
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
