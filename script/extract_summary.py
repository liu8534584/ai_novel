import os
import re
import json
import requests
from typing import List, Dict

# --- 配置区 ---
CONFIG = {
    "input_file": "/Users/liuda/Downloads/2.txt",        # 你的小说原始文件
    "output_file": "/Users/liuda/Downloads/train_novel.jsonl",  # 输出的训练数据集
    "api_url": "http://localhost:1234/v1/chat/completions", # LM Studio 默认端口
    "max_content_len": 4000,             # 喂给模型提取摘要的最大字符数
    "min_chapter_len": 300               # 忽略少于此字数的章节（通常是公告）
}

def read_file_safely(file_path):
    # 优先尝试 UTF-8
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()
    except UnicodeDecodeError:
        # 如果报错，尝试中文 GB18030 编码
        try:
            with open(file_path, 'r', encoding='gb18030', errors='ignore') as f:
                return f.read()
        except Exception as e:
            print(f"文件读取失败: {e}")
            return None

def split_chapters(file_path: str) -> List[Dict]:
    """使用正则拆分章节，返回标题和内容的列表"""
    content = read_file_safely(file_path)

    # 正则：匹配 第xx章、第xx节、或者数字开头带、的标题
    # 这个正则兼容性较强，可根据你的小说格式微调
    pattern = r'(第[一二三四五六七八九十百千万\d]+[章节回].*?)\n'

    parts = re.split(pattern, content)
    chapters = []

    # split 会将匹配项和非匹配项交替排列
    # parts[0] 通常是序言，从 parts[1] 开始是 标题, 内容, 标题, 内容...
    for i in range(1, len(parts), 2):
        title = parts[i].strip()
        body = parts[i+1].strip() if i+1 < len(parts) else ""
        if len(body) > CONFIG["min_chapter_len"]:
            chapters.append({"title": title, "content": body})

    return chapters

def get_llm_summary(title: str, content: str) -> str:
    truncated_content = content[:CONFIG["max_content_len"]]

    # 更加严苛的指令，防止模型输出废话
    prompt = f"""阅读以下小说章节，直接输出摘要和关键词。
要求：
1. 摘要：150字以内，直接陈述剧情（人物、冲突、结果）。
2. 关键词：3个核心词，用竖线 | 分隔。
3. 禁止输出 <think> 标签，禁止输出任何开场白。

【章节标题】：{title}
【正文片段】：{truncated_content}
"""

    payload = {
        "messages": [
            {"role": "system", "content": "你是一个只输出结构化摘要的机器人，不准聊天，不准思考。"},
            {"role": "user", "content": prompt}
        ],
        "temperature": 0.1,
    }

    try:
        response = requests.post(CONFIG["api_url"], json=payload, timeout=300)
        full_text = response.json()['choices'][0]['message']['content']

        # --- 核心清理逻辑 ---
        # 1. 剔除 <think> 标签及其内容
        clean_text = re.sub(r'<think>.*?</think>', '', full_text, flags=re.DOTALL)
        # 2. 剔除可能的 Markdown 代码块标记
        clean_text = clean_text.replace("```", "").strip()

        return clean_text
    except Exception as e:
        print(f"提取章节《{title}》失败: {e}")
        return None

def main():
    print(f"开始读取: {CONFIG['input_file']}...")
    chapters = split_chapters(CONFIG["input_file"])
    print(f"成功拆分出 {len(chapters)} 个章节。")

    with open(CONFIG["output_file"], 'w', encoding='utf-8') as f_out:
        for idx, ch in enumerate(chapters):
            print(f"[{idx+1}/{len(chapters)}] 正在提炼: {ch['title']}")

            summary = get_llm_summary(ch['title'], ch['content'])

            if summary:
                # 构造标准微调格式 (Qwen/Llama 通用)
                training_point = {
                    "messages": [
                        {"role": "system", "content": "你是一位专业的网络小说作家，请根据标题和摘要续写正文。"},
                        {"role": "user", "content": f"【标题】：{ch['title']}\n【剧情提要】：{summary}"},
                        {"role": "assistant", "content": ch['content']}
                    ]
                }
                f_out.write(json.dumps(training_point, ensure_ascii=False) + "\n")
                f_out.flush() # 实时保存，防止奔溃

    print(f"所有任务已完成！输出文件: {CONFIG['output_file']}")

if __name__ == "__main__":
    main()