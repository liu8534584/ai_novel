import json
import re
import os

def clean_text(text):
    """剔除小说中的广告和乱码"""
    # 过滤常见的网文广告词
    ad_patterns = [
        r"求月票", r"求收藏", r"加群", r"http[s]?://\S+",
        r"点击下一页", r"本章未完", r"（.*）"
    ]
    for pattern in ad_patterns:
        text = re.sub(pattern, "", text)

    # 压缩多余空行和空格
    text = re.sub(r'\n+', '\n', text)
    text = re.sub(r' +', ' ', text)
    return text.strip()

def process_novels(input_dir, output_file, seq_length=1500):
    dataset = []

    for filename in os.listdir(input_dir):
        if filename.endswith(".txt"):
            file_path = os.path.join(input_dir, filename)
            print(f"正在处理: {filename}")

            # --- 兼容性读取逻辑开始 ---
            content = ""
            try:
                # 首先尝试 UTF-8
                with open(file_path, 'r', encoding='utf-8') as f:
                    content = f.read()
            except UnicodeDecodeError:
                # 如果 UTF-8 失败，尝试 GB18030 (兼容 GBK)
                try:
                    with open(file_path, 'r', encoding='gb18030') as f:
                        content = f.read()
                except Exception as e:
                    print(f"警告：跳过文件 {filename}，无法识别编码。错误: {e}")
                    continue
            # --- 兼容性读取逻辑结束 ---

            content = clean_text(content)

            # 按固定长度切片
            for i in range(0, len(content), seq_length):
                chunk = content[i : i + seq_length + 500]
                last_punc = max(chunk.rfind('。'), chunk.rfind('！'), chunk.rfind('？'))
                if last_punc == -1: continue

                actual_chunk = chunk[:last_punc + 1]

                data_point = {
                    "text": f"【末世小说演练】\n【书名】：{filename.replace('.txt','')}\n【正文】：\n{actual_chunk}"
                }
                dataset.append(data_point)

    # 写入 JSONL
    with open(output_file, 'w', encoding='utf-8') as f:
        for entry in dataset:
            f.write(json.dumps(entry, ensure_ascii=False) + '\n')

    print(f"处理完成！共生成 {len(dataset)} 条逻辑片段。")

# 使用方法
if __name__ == "__main__":
    # 请确保你的50本小说放在 raw_novels 文件夹里
    process_novels("/Users/liuda/Documents/txt/末世文_all", "/Users/liuda/Documents/txt/末世文_all/train.jsonl")