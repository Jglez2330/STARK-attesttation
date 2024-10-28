import re

import angr
from angr.knowledge_plugins.cfg import CFGNode
import matplotlib.pyplot as plt
import networkx as nx
def main():
    proj = angr.Project('/Users/jglez2330/Library/Mobile Documents/com~apple~CloudDocs/personal/STARK-attesttation/test/a.out', load_options={'auto_load_libs': False})
    cfg = proj.analyses.CFGFast(normalize=True, show_progressbar=True, resolve_indirect_jumps=True)

    #end_state, cfg, path = get_execution_path(proj, cfg) # get example execution path
    #print("Execution path ended because it %s" % end_state)

    graph = cfg.graph
    adj = nx.generate_adjlist(graph)
    adj = adjlist_to_dict(adj)
    print(adj)
    # Visualize the CFG using NetworkX and Matplotlib
    plt.figure(figsize=(12, 8))
    nx.draw(graph, with_labels=True, node_size=500, font_size=8, node_color='lightblue', edge_color='gray')
    plt.title("Control-Flow Graph (CFG) Extracted with angr")
    plt.show()



if __name__ == "__main__":
    main()