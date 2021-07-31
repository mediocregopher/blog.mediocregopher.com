---
layout: default
---

{%- if site.posts.size > 0 -%}
<ul id="posts-list">
  {%- for post in site.posts -%}
  {%- if post.hide -%}{%- continue -%}{%- endif -%}
  <li>
    <h2>
      <a href="{{ post.url | relative_url }}">
        {{ post.title | escape }}
      </a>
    </h2>
    <span>{{ post.date | date: site.date_format }}</span>
    {%- if post.updated %}
    <span>(Updated {{ post.updated | date: site.date_format }})</span>
    {% endif -%}
    <p>{{ post.description }}</p>
  </li>
  {%- endfor -%}
</ul>
{%- endif -%}
