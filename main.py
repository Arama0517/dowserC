import json
import re
import subprocess
import sys
from pathlib import Path
from platform import system
from time import sleep
from typing import NoReturn


def pause_and_exit(code) -> NoReturn:
    input('按回车键继续...: ')
    sys.exit(code)


# 错误代码
NOT_SUPPORTED_SYSTEM = 1
LAUNCHER_SETTINGS_PATH_NOT_FOUND = 2

# 检查系统
if system() != 'Windows':
    print('此程序仅支持Windows')
    pause_and_exit(NOT_SUPPORTED_SYSTEM)

# 定义路径
CWD = Path.cwd()
DATA_DIR = CWD.joinpath('Paradox Interactive', 'Hearts of Iron IV')
MOD_DIR = CWD.joinpath('mod')

# 创建目录
DATA_DIR.mkdir(parents=True, exist_ok=True)
MOD_DIR.mkdir(parents=True, exist_ok=True)

# 检查启动器配置文件
LAUNCHER_SETTINGS_PATH = CWD.joinpath('launcher-settings.json')

if not LAUNCHER_SETTINGS_PATH.exists():
    print('启动器配置文件不存在')
    print('请确认dowserC.exe是否在游戏根目录下')
    pause_and_exit(LAUNCHER_SETTINGS_PATH_NOT_FOUND)

# 修改启动器配置文件
with LAUNCHER_SETTINGS_PATH.open('r', encoding='utf-8') as f:
    content = json.load(f)
with LAUNCHER_SETTINGS_PATH.open('w', encoding='utf-8') as f:
    content['gameDataPath'] = str(DATA_DIR)
    json.dump(content, f, indent=4, ensure_ascii=False)

# 检测无效的定位文件
for mod in DATA_DIR.joinpath('mod').iterdir():
    if mod.is_file() and mod.suffix == '.mod':
        with mod.open('r', encoding='utf-8') as f:
            content = f.read()
            path = re.search(r'path\s*=\s*"([^"]+)"', content)
            if path:
                path = path.group(1)
            else:
                print(f'检测到无效的定位文件: {mod.name}')
                mod.unlink()
                continue

        if not Path(path).exists():
            print(f'检测到无效的定位文件: {mod.name}')
            mod.unlink()
            continue

# 加载Mod
mod_names = []
for mod in MOD_DIR.iterdir():
    descriptor_path = mod.joinpath('descriptor.mod')
    if mod.is_dir() and descriptor_path.exists():
        with descriptor_path.open('r', encoding='utf-8') as f:
            content = f.read()
            path = re.search(r'name\s*=\s*"([^"]+)"', content)
            if path:
                path = path.group(1)
                mod_names.append(path)
            else:
                continue

        with DATA_DIR.joinpath('mod', f'{path}.mod').open('w', encoding='utf-8') as f:
            f.write(content)
            f.write(f'\npath="{mod.as_posix()}"')


# 输出信息
print(f'共计加载: {len(mod_names)}个模组')
print('他们分别是: ')
print('-' * 50)
for path in mod_names:
    print(path)
print('-' * 50)

# 启动dowser(启动器)
print('准备启动游戏!3秒后自动退出...')
subprocess.run(CWD.joinpath('dowser.exe'), check=False)
sleep(3)
sys.exit(0)